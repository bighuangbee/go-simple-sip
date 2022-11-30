package sip

import (
	"github.com/bighuangbee/go-simple-sip/internal/sip/message"
	"net"
)

//IPC User Agent Client
//UDP连接信息 或 TCP客户端句柄的索引
type UacConn net.UDPAddr

type Transport interface {
	Write(uacMsg *UacMsg) error
}

type Gb28181Serve interface {
	Catalog(uacMsg *UacMsg, catalog *message.Query) error
	Play(uacMsg *UacMsg, req *message.PlayReq) (streamId string, err error)
	PlayRespone(uacMsg *UacMsg) (err error)
}
