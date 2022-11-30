package sip

import (
	"github.com/jart/gosip/sip"
	"net"
)

//User Agent Client 请求报文
type UacRequest struct {
	uac     net.Addr //IPC连接地址
	message []byte
}

//User Agent Client SIP信令消息
type UacMsg struct {
	uacAddr net.Addr //IPC连接地址,User Agent
	msg     *sip.Msg
}

func NewUacMsg(uacConn net.Addr, msg *sip.Msg) *UacMsg {
	return &UacMsg{uacAddr: uacConn, msg: msg}
}

func (this *UacRequest) ToUacMsg() (*UacMsg, error) {
	sipMsg, err := sip.ParseMsg(this.message)
	if err != nil {
		return nil, err
	}
	return &UacMsg{uacAddr: this.uac, msg: sipMsg}, nil
}
