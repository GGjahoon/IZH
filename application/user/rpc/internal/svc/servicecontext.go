package svc

import (
	"github.com/GGjahoon/IZH/application/user/rpc/internal/config"
	"github.com/GGjahoon/IZH/application/user/rpc/internal/model"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config, userModel model.UserModel) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		UserModel: userModel,
	}
}
