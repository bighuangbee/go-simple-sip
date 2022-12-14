package sip

import (
	"fmt"
	"github.com/bighuangbee/go-simple-sip/internal/sip/message"
	"github.com/bighuangbee/go-simple-sip/pkg/tools"
	"github.com/jart/gosip/sdp"
	"github.com/jart/gosip/sip"
	"net"
)

func (this *Uas) Play(uacAddr net.Addr, req *message.PlayReq) (streamId string, err error) {
	addr, err := net.ResolveUDPAddr(uacAddr.Network(), uacAddr.String())
	if err != nil{
		return "", err
	}

	ssrc := this.GenSSRC(1)

	playSdp := sdp.New(addr)
	playSdp.Origin = sdp.Origin{
		User:    this.Bootstrap.Server.GB28181.SipId,
		Addr:    this.Bootstrap.Server.Media.Addr,
		ID:      "0",
		Version: "0",
	}
	playSdp.Addr = this.Bootstrap.Server.Media.Addr
	playSdp.Audio = nil
	playSdp.Video = &sdp.Media{
		Proto: "TCP/RTP/AVP",
		Port:  this.Bootstrap.Server.Media.StreamRecvPort,
		Codecs: []sdp.Codec{
			{PT: 96, Name: "PS", Rate: 90000},
			{PT: 98, Name: "H264", Rate: 90000},
			{PT: 97, Name: "MPEG4", Rate: 90000},
		},
	}
	playSdp.Session = "Play"
	playSdp.Time = "0 0"
	playSdp.SendOnly = false
	playSdp.RecvOnly = true
	playSdp.Attrs = [][2]string{[2]string{"setup", "passive"}, [2]string{"connection", "new"}}
	playSdp.Other = [][2]string{[2]string{"y", ssrc}}

	sipPlay := new(sip.Msg)
	sipPlay.CallID = tools.Rand(32)
	sipPlay.CSeq = 12
	sipPlay.Request = &sip.URI{
		User: req.ChannelId,
		Host: this.Bootstrap.Server.Media.Addr,
	}

	sipPlay.Subject = fmt.Sprintf("%s:%s,%s:%s", req.ChannelId, ssrc, this.Bootstrap.Server.GB28181.SipId, ssrc)
	//sipPlay.Via = uacMsg.msg.Via //branch事务ID
	sipPlay.Via = &sip.Via{
		Protocol:  "SIP",
		Version:   "2.0",
		Transport: "UDP",
		Host:      req.Addr,
		Port:      req.Port,
		Param: &sip.Param{
			Name:  "branch",
			Value: message.GenBranch(),
			Next:  &sip.Param{Name: "rport"},
		},
	}

	sipPlay.Payload = &sip.MiscPayload{
		T: message.SDP,
		D: playSdp.Data(),
	}

	sipPlay.Method = "INVITE"
	sipPlay.CSeqMethod = "INVITE"
	sipPlay.To = &sip.Addr{
		Uri: &sip.URI{
			User: req.ChannelId,
			Host: req.Addr,
			Port: req.Port,
		},
	}
	sipPlay.From = &sip.Addr{
		Uri: &sip.URI{
			Scheme: "sip",
			User:   this.Bootstrap.Server.GB28181.SipId,
			Host:   this.Bootstrap.Server.GB28181.SipDomain,
		},
		Param: &sip.Param{
			Name:  "tag",
			Value: tools.Rand(32),
		},
	}
	sipPlay.Contact = sipPlay.From

	return message.SsrcTostreamId(ssrc), this.Transport.Write(&UacMsg{
		uacAddr: uacAddr,
		msg:     sipPlay,
	})

}

func (this *Uas) PlayRespone(uacMsg *UacMsg) (err error) {
	m := uacMsg.msg.Copy()
	m.Request = m.From.Uri
	m.Status = 0
	m.Method = "ACK"
	m.CSeqMethod = "ACK"
	m.Payload = nil
	m.From.Uri.User = this.Bootstrap.Server.GB28181.SipId
	m.Via.Port = this.Bootstrap.Server.GB28181.Port

	return this.Transport.Write(&UacMsg{
		uacAddr: uacMsg.uacAddr,
		msg:     m,
	})
}
