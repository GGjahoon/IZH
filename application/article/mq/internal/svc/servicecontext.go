package svc

import (
	"github.com/GGjahoon/IZH/application/article/mq/internal/config"
	"github.com/GGjahoon/IZH/application/article/mq/internal/model"
	"github.com/GGjahoon/IZH/application/user/rpc/user"
	"github.com/GGjahoon/IZH/pkg/es"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config       config.Config
	ArticleModel model.ArticleModel
	BizRedis     *redis.Redis
	UserRPC      user.User
	Es           *es.Es
}

func NewServiceContext(config config.Config,
	articleModel model.ArticleModel,
	redis *redis.Redis,
	userRPC user.User,
	es *es.Es,
) *ServiceContext {
	return &ServiceContext{
		Config:       config,
		ArticleModel: articleModel,
		BizRedis:     redis,
		UserRPC:      userRPC,
		Es:           es,
	}
}
