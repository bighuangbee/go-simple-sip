package sip

import (
	"context"
	"errors"
	"fmt"
	domain2 "github.com/bighuangbee/go-simple-sip/internal/data/domain"
	"github.com/bighuangbee/go-simple-sip/internal/sip/message"
	"github.com/jart/gosip/sip"
)

//向UAC发送catalog请求
func (this *Uas) Catalog(uacMsg *UacMsg, catalog *message.Query) error {

	queryCatalog := uacMsg.msg.Copy()
	queryCatalog.Method = sip.MethodMessage
	queryCatalog.CSeqMethod = sip.MethodMessage
	queryCatalog.Via.Port = queryCatalog.From.Uri.Port
	queryCatalog.Status = 0
	queryCatalog.From.Uri.User = this.SysConf.GB28181.SipId
	queryCatalog.From.Uri.Host = this.SysConf.GB28181.SipDomain
	queryCatalog.From.Uri.Port = 0
	queryCatalog.To = uacMsg.msg.From
	queryCatalog.To.Param = nil
	queryCatalog.Payload = &sip.MiscPayload{
		T: message.MANSCDP,
		D: message.Marshal(catalog),
	}

	uacMsg.msg = queryCatalog
	if err := this.Transport.Write(uacMsg); err != nil {
		return errors.New("Catalog " + err.Error())
	}
	return nil
}

func (this *Uas) CatalogRespone(uacMsg *UacMsg) error {
	payload := uacMsg.msg.Payload.Data()
	catalogRespone := &message.CatalogResponse{}
	message.Unmarshal(payload, catalogRespone)

	if len(catalogRespone.DeviceList.Channels) > 0 {
		c := catalogRespone.DeviceList.Channels[0]
		device := domain2.Device{
			DeviceId:     catalogRespone.DeviceID,
			Name:         c.Name,
			Manufacturer: c.Manufacturer,
			Model:        c.Model,
			//Firmware:      "",
			//Transport:     "",
			Status:      message.StatusMap(c.Status),
			HostAddress: c.IPAddress,
			Ip:          c.IPAddress,
			Port:        uint16(uacMsg.uacConn.Port),
			//Expires:       0,
			//Charset:       "",
		}
		this.Repo.Device.Save(context.Background(), &device)

		for _, channle := range catalogRespone.DeviceList.Channels {
			fmt.Println("0------------ ", channle.IPAddress, uacMsg.uacConn.Port)

			c := domain2.Channel{
				DeviceId:     catalogRespone.DeviceID,
				ChannelId:    channle.DeviceID,
				Name:         c.Name,
				Manufacturer: c.Manufacturer,
				Model:        c.Model,
				//Firmware:      "",
				//Transport:     "",
				Status:      message.StatusMap(c.Status),
				HostAddress: fmt.Sprintf("%s:%d", uacMsg.uacConn.IP.String(), uacMsg.uacConn.Port),
				Ip:          uacMsg.uacConn.IP.String(),
				Port:        uint16(uacMsg.uacConn.Port),
				//Expires:       0,
				//Charset:       "",
			}
			this.Repo.Channel.Save(context.Background(), &c)

		}
	}

	//回复200
	respone := new(sip.Msg)
	respone.Via = uacMsg.msg.Via
	respone.Status = sip.StatusOK
	respone.CSeq = uacMsg.msg.CSeq
	respone.CSeqMethod = sip.MethodMessage
	respone.CallID = uacMsg.msg.CallID
	respone.From = uacMsg.msg.From
	respone.To = uacMsg.msg.To
	return this.Transport.Write(&UacMsg{uacMsg.uacConn, respone})
}
