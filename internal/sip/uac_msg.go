package sip

import (
	"github.com/jart/gosip/sip"
)

//User Agent Client 请求报文
type UacRequest struct {
	uac     *UacConn //IPC连接地址
	message []byte
}

//User Agent Client SIP信令消息
type UacMsg struct {
	uacConn *UacConn //IPC连接地址,User Agent
	msg     *sip.Msg
}

func NewUacMsg(uacConn *UacConn, msg *sip.Msg) *UacMsg {
	return &UacMsg{uacConn: uacConn, msg: msg}
}

func (this *UacRequest) ToUacMsg() (*UacMsg, error) {
	sipMsg, err := sip.ParseMsg(this.message)
	if err != nil {
		return nil, err
	}
	return &UacMsg{uacConn: this.uac, msg: sipMsg}, nil
}
