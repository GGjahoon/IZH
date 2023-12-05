package svc

import (
	"github.com/GGjahoon/IZH/application/article/rpc/internal/config"
	"github.com/GGjahoon/IZH/application/article/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"golang.org/x/sync/singleflight"
)

type ServiceContext struct {
	Config            config.Config
	ArticleModel      model.ArticleModel
	BizRedis          *redis.Redis
	SingleFlightGroup singleflight.Group
}

func NewServiceContext(c config.Config, articlModel model.ArticleModel, rds *redis.Redis) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		ArticleModel: articlModel,
		BizRedis:     rds,
	}
}
