package logic_test

import (
	"testing"

	"github.com/GGjahoon/IZH/application/user/rpc/internal/config"
	"github.com/GGjahoon/IZH/application/user/rpc/internal/model"
	"github.com/GGjahoon/IZH/application/user/rpc/internal/server"
	"github.com/GGjahoon/IZH/application/user/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func NewTestServer(t *testing.T, userModel model.UserModel) *server.UserServer {
	c := config.Config{
		BizRedis: redis.RedisConf{
			Host: "127.0.0.1:6379",
			Pass: "",
			Type: "node",
		},
	}
	ctx := svc.NewServiceContext(c, userModel)
	testSever := server.NewUserServer(ctx)
	return testSever
}
