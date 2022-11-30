package model

import (
	domain2 "github.com/bighuangbee/go-simple-sip/internal/data/domain"
)

type Repo struct {
	Device  domain2.IDeviceRepo
	Channel domain2.IChannelRepo
	Stream  domain2.IStreamsRepo
}
