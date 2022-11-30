package service

import (
	"fmt"
	"github.com/bighuangbee/go-simple-sip/internal/config"
	"github.com/bighuangbee/go-simple-sip/internal/data"
	model2 "github.com/bighuangbee/go-simple-sip/internal/data/model"
	"github.com/bighuangbee/go-simple-sip/pkg/tools/log"
)

const messagePoolSize = 1000

type Service struct {
	//系统配置
	SysConf *config.SysConf
	//rtp SSRC
	SSRC int
	//数据层
	data *data.Data
	//repo
	Repo *model2.Repo
}

func New(sysConf *config.SysConf) *Service {
	sipServer := &Service{SysConf: sysConf} //messagePool: make(chan *UacRequest, messagePoolSize),

	sipServer.SSRC = 10

	logger := log.New("./logs")
	db, err := data.New(&data.Options{
		Address:  sysConf.Database.Address,
		UserName: sysConf.Database.UserName,
		Password: sysConf.Database.Password,
		DBName:   sysConf.Database.DBName,
		Logger:   logger,
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	sipServer.data = &data.Data{
		Db: db,
	}
	sipServer.Repo = &model2.Repo{
		Device:  &model2.DeviceRepo{Data: sipServer.data},
		Channel: &model2.ChannelRepo{Data: sipServer.data},
		Stream:  &model2.StreamRepo{Data: sipServer.data},
	}
	return sipServer
}
