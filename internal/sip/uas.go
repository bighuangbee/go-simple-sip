package sip

import (
	"context"
	"errors"
	"fmt"
	"github.com/bighuangbee/go-simple-sip/internal/config"
	"github.com/bighuangbee/go-simple-sip/internal/service"
	"github.com/bighuangbee/go-simple-sip/internal/sip/message"
	"github.com/bighuangbee/go-simple-sip/pkg/tools"
	"github.com/jart/gosip/sdp"
	"github.com/jart/gosip/sip"
	"golang.org/x/sync/errgroup"
	"runtime"
	"strconv"
)


type Uas struct {
	*service.Service
	Transport Transport
	uacMsgQueue chan *UacMsg
}

func NewUas(sysConf *config.SysConf) *Uas {
	uacMsgQueue := make(chan *UacMsg, runtime.NumGoroutine() * 100)
	return &Uas{
		Transport: NewUdpTransport(sysConf.Server.UpdPort, uacMsgQueue),
		Service: service.New(sysConf),
		uacMsgQueue: uacMsgQueue,
	}
}

//isRTP 0=实时流，1=历史流
func (this *Uas) GenSSRC(isRTP int) string {
	this.SSRC += 1
	return strconv.Itoa(isRTP) + tools.RepairSuff(this.SSRC, 9)

}


func (this *Uas) Run() {
	eg := errgroup.Group{}
	for i := 0; i < 1; i++ {
		eg.Go(func() error {
			for {
				select {
				case uacMsg := <-this.uacMsgQueue:
					{
						if err := this.uacMsgHandler(uacMsg); err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		})
	}
	eg.Wait()
}

func (this *Uas) uacMsgHandler(uacMsg *UacMsg) (err error) {

	if uacMsg.msg.Method == sip.MethodRegister {
		err = this.Register(uacMsg)

		//time.AfterFunc(time.Second, func() {
		//	streamId, _ := this.Play(uacMsg.uacConn, &gb.PlayReq{
		//		ChannelId:  uacMsg.msg.From.Uri.User,
		//		Addr:      uacMsg.msg.Request.Host,
		//		Port:      uacMsg.msg.Request.Port,
		//	})
		//	fmt.Println("------------------streamId:",streamId)
		//})
	} else {
		if uacMsg.msg.Payload != nil {
			payload := uacMsg.msg.Payload.Data()

			if uacMsg.msg.Status == sip.StatusOK {
				var sdpMsg *sdp.SDP
				if sdpMsg, err = sdp.Parse(string(payload)); err != nil {
					return err
				}

				if sdpMsg.Session == "Play" {
					err = this.PlayRespone(uacMsg)
				}
			} else {
				base := message.Payload{}
				if err = message.Unmarshal(payload, &base); err != nil {
					return errors.New("CmdType parse " + err.Error())
				}

				switch base.CmdType {
				case message.CmdTypeKeepalive:
					{
						err = this.Keepalive(uacMsg)
						if err != nil {
							return err
						}

						deviceId := uacMsg.msg.From.Uri.User
						_, total, err := this.Service.Repo.Channel.List(context.Background(), deviceId)
						fmt.Println("-------------- ", err, total, deviceId)
						if err != nil {
							return err
						}
						if total == 0 {
							err = this.Catalog(uacMsg, &message.Query{
								Payload: message.Payload{
									CmdType:  message.CmdTypeCatalog,
									SN:       base.SN,
									DeviceID: deviceId,
								},
							})
						}

						//todo 认证、存储
					}
				case message.CmdTypeCatalog:
					err = this.CatalogRespone(uacMsg)

				}
			}

		}
	}

	return
}



func (this *Uas) Keepalive(uacMsg *UacMsg) error {
	uacMsg.msg.Status = sip.StatusOK
	uacMsg.msg.Payload = nil

	if err := this.Transport.Write(uacMsg); err != nil {
		return errors.New("Keepalive " + err.Error())
	}
	return nil
}

func (this *Uas) Register(uacMsg *UacMsg) error {
	uacMsg.msg.Status = sip.StatusOK
	//回复
	if err := this.Transport.Write(uacMsg); err != nil {
		return errors.New("Register " + err.Error())
	}
	return nil
}
