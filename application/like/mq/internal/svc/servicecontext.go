package svc

import "github.com/GGjahoon/IZH/application/like/mq/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(config config.Config) *ServiceContext {
	return &ServiceContext{
		Config: config,
	}
}
