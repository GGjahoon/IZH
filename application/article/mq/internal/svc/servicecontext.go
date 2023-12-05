package svc

import (
	"github.com/GGjahoon/IZH/application/article/mq/internal/config"
	"github.com/GGjahoon/IZH/application/article/mq/internal/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config       config.Config
	ArticleModel model.ArticleModel
	BizRedis     *redis.Redis
}

func NewServiceContext(config config.Config, articleModel model.ArticleModel, redis *redis.Redis) *ServiceContext {
	return &ServiceContext{
		Config:       config,
		ArticleModel: articleModel,
		BizRedis:     redis,
	}
}
