package sip

import (
	"github.com/bighuangbee/go-simple-sip/internal/sip/message"
	"net"
)


type Transport interface {
	Write(uacMsg *UacMsg) error
}

type Gb28181Serve interface {
	Catalog(uacMsg *UacMsg, catalog *message.Query) error
	Play(uacAddr net.Addr, req *message.PlayReq) (streamId string, err error)
	PlayRespone(uacMsg *UacMsg) (err error)
}
