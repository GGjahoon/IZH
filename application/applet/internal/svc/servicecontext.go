package svc

import (
	"github.com/GGjahoon/IZH/application/applet/internal/config"
	"github.com/GGjahoon/IZH/application/user/rpc/user"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	UserRpc  user.User
	BizReids *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		UserRpc:  user.NewUser(zrpc.MustNewClient(c.UserRPC)),
		BizReids: redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
	}
}
