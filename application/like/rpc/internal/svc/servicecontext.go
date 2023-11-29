package svc

import (
	"github.com/GGjahoon/IZH/application/like/rpc/internal/config"
	"github.com/zeromicro/go-queue/kq"
)

type ServiceContext struct {
	Config         config.Config
	KqPusherClient *kq.Pusher
}

func NewServiceContext(c config.Config, kqPusherClient *kq.Pusher) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		KqPusherClient: kqPusherClient,
	}
}
