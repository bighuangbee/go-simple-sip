package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bighuangbee/go-simple-sip/internal/config"
	domain2 "github.com/bighuangbee/go-simple-sip/internal/data/domain"
	"github.com/bighuangbee/go-simple-sip/internal/sip"
	"github.com/bighuangbee/go-simple-sip/internal/sip/message"
	"github.com/bighuangbee/go-simple-sip/pkg/tools"
	"github.com/gin-gonic/gin"
	"net"
)

var sipUas *sip.Uas

func main() {
	var flagConf string
	flag.StringVar(&flagConf, "conf", "../conf/config.yaml", "config path, eg: -conf config.yaml")
	flag.Parse()

	var bc config.Bootstrap
	if err := tools.InitConfig(flagConf, &bc); err != nil{
		fmt.Println("InitConfig err", err)
		return
	}

	go httpServer(bc.Server.HttpPort)

	sipUas = sip.NewUas(&bc)
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
		c.JSON(0, "")
		return
	}

	if stream.StreamId == "" {
		channel, err := sipUas.Repo.Channel.GetByDeviceId(context.Background(), deviceId, channelId)
		if err != nil {
			fmt.Println(err)
			c.JSON(0, "")
			return
		}

		fmt.Println("channel ,", channel)


		addr, err := net.ResolveUDPAddr("udp", channel.HostAddress)
		if err != nil{
			fmt.Println("net.ResolveUDPAddr", err, channel.HostAddress)
			return
		}

		streamId, err := sipUas.Play(addr, &message.PlayReq{
			DeviceId:  deviceId,
			ChannelId: channelId,
			Addr:      channel.Ip,
			Port:      channel.Port,
		})
		if err != nil {
			fmt.Println("uasServer.Play ", err)
			return
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
				Hls:        fmt.Sprintf("http://%s:%d/rtp/%s/hls.m3u8", sipUas.Bootstrap.Server.Media.Addr, sipUas.Bootstrap.Server.Media.Port, streamId),
			Rtmp:       "",
			Rtsp:       "",
			Wsflv:      fmt.Sprintf("ws://%s:%d/rtp/%s.live.flv", sipUas.Bootstrap.Server.Media.Addr, sipUas.Bootstrap.Server.Media.Port, streamId),
			Stream:     0,
		}

		sipUas.Repo.Stream.Create(context.Background(), stream)
	}

	c.JSON(0, stream)
}
