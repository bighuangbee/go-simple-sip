package sip

import (
	"fmt"
	"github.com/jart/gosip/sip"
	"net"
)

type UdpTransport struct {
	udpConn *net.UDPConn
}

const bufferSize uint16 = 65535 - 20 - 8 // IPv4 max size - IPv4 Header size - UDP Header size

func NewUdpTransport(port uint16, uacMsgQueue chan *UacMsg) *UdpTransport {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			var buf = make([]byte, bufferSize)
			n, uac, err := udpConn.ReadFromUDP(buf[0:])
			if err != nil {
				fmt.Println("read udp failed, err", err)
				continue
			}

			fmt.Println(fmt.Sprintf("【Reqeust】form: %s, data:\n%s\n", uac.String(), string(buf[:n])))

			msg, err := sip.ParseMsg(buf[:n])
			if err != nil {
				fmt.Println("sip.ParseMsg err:", err)
				continue
			}

			uacMsgQueue <- &UacMsg{
				uacAddr: uac,
				msg:     msg,
			}
		}
	}()

	fmt.Println("【Sip Uas】Run Success, Listening Port:", port)

	return &UdpTransport{udpConn: udpConn}
}

func (this *UdpTransport) Write(uacMsg *UacMsg) error {
	fmt.Println(fmt.Sprintf("-------------Write:\n%s\n------Write end\n", uacMsg.msg.String()))
	_, err := this.udpConn.WriteTo([]byte(uacMsg.msg.String()), uacMsg.uacAddr)
	return err
}
