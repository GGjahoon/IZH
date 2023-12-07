package svc

import (
	"github.com/GGjahoon/IZH/application/follow/rpc/internal/config"
	"github.com/GGjahoon/IZH/application/follow/rpc/internal/model"
	"github.com/GGjahoon/IZH/pkg/orm"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config           config.Config
	DB               *orm.DB
	FollowModel      *model.FollowModel
	FollowCountModel *model.FollowCountModel
	BizRedis         *redis.Redis
}

func NewServiceContext(c config.Config,
	db *orm.DB,
	followModel *model.FollowModel,
	followCountModel *model.FollowCountModel,
	rds *redis.Redis,
) *ServiceContext {
	return &ServiceContext{
		Config:           c,
		DB:               db,
		FollowModel:      followModel,
		FollowCountModel: followCountModel,
		BizRedis:         rds,
	}
}
