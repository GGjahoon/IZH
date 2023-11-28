package svc

import (
	"github.com/GGjahoon/IZH/application/article/rpc/internal/config"
	"github.com/GGjahoon/IZH/application/article/rpc/internal/model"
)

type ServiceContext struct {
	Config       config.Config
	ArticleModel model.ArticleModel
}

func NewServiceContext(c config.Config, articlModel model.ArticleModel) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		ArticleModel: articlModel,
	}
}
