package main

import (
	"context"
	"fmt"
	"github.com/bighuangbee/go-simple-sip/internal/config"
	domain2 "github.com/bighuangbee/go-simple-sip/internal/data/domain"
	"github.com/bighuangbee/go-simple-sip/internal/sip"
	"github.com/bighuangbee/go-simple-sip/internal/sip/message"
	"github.com/gin-gonic/gin"
	"net"
)

var sipUas *sip.Uas

func main() {
	sysConf := config.SysConf{
		Server: &config.Server{
			UpdPort:  5060,
			HttpPort: 6660,
			HttpAddr: "192.168.80.242", //

		},
		GB28181: &config.GB28181{
			SipId:     "37070000082008000001",
			SipDomain: "3707000008",
		},
		Media: &config.Media{
			Addr:           "192.168.80.242",
			Port:           8080,
			StreamRecvPort: 10000,
		},
		Database: &config.Database{
			Address:  "localhost",
			UserName: "root",
			Password: "Hiscene2022",
			DBName:   "gbsip",
			Driver:   "mysql",
			Timeout:  10,
		},
	}

	go httpServer(sysConf.Server.HttpPort)

	sipUas = sip.NewUas(&sysConf)
	sipUas.Run()
}

func httpServer(port uint16) {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, 111")
	})
	r.GET("/channels/:deviceId", channnels)
	r.GET("/channels/stream/:deviceId/:channelId", streams)
	err := r.Run(":" + fmt.Sprintf("%d", port))
	fmt.Println("httpServer: ", err)
}

func channnels(c *gin.Context) {
	deviceId := c.Param("deviceId")
	fmt.Println("-----deviceId:", deviceId)
}

func streams(c *gin.Context) {
	deviceId := c.Param("deviceId")
	channelId := c.Param("channelId")
	fmt.Println("-----deviceId:", deviceId, channelId)

	stream, err := sipUas.Repo.Stream.GetByDeviceId(context.Background(), &domain2.ChannelQuery{
		DeviceId:  deviceId,
		ChannelId: channelId,
	})
	if err != nil {
		fmt.Println("GetByDeviceId:", err)
		return
	}

	if stream.StreamId == "" {
		channel, err := sipUas.Repo.Channel.GetByDeviceId(context.Background(), deviceId, channelId)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("channel ,", channel)




		streamId, err := sipUas.Play(sip.NewUacMsg(&sip.UacConn{
			IP:   net.ParseIP(channel.Ip).To4(),
			Port: int(channel.Port),
		}, nil), &message.PlayReq{
			DeviceId:  deviceId,
			ChannelId: channelId,
			Addr:      channel.Ip,
			Port:      channel.Port,
		})
		if err != nil {
			fmt.Println("uasServer.Play ", err)
		}

		fmt.Println("--------------------- streamId", streamId)

		stream = &domain2.Stream{
			T:          0,
			DeviceId:   deviceId,
			ChannelId:  channelId,
			StreamType: "",
			Status:     0,
			Callid:     "",
			Stop:       0,
			Msg:        "",
			Cseqno:     0,
			StreamId:   streamId,
			Hls:        fmt.Sprintf("http://%s:%d/rtp/%s/hls.m3u8", sipUas.SysConf.Media.Addr, sipUas.SysConf.Media.Port, streamId),
			Rtmp:       "",
			Rtsp:       "",
			Wsflv:      fmt.Sprintf("ws://%s:%d/rtp/%s.live.flv", sipUas.SysConf.Media.Addr, sipUas.SysConf.Media.Port, streamId),
			Stream:     0,
		}

		sipUas.Repo.Stream.Create(context.Background(), stream)
	}

	c.JSON(0, stream)
}
