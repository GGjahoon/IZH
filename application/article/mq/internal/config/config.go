package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	service.ServiceConf
	KqConsumerConf        kq.KqConf
	ArticleKqConsumerConf kq.KqConf
	DataSource            string
	BizRedisConf          redis.RedisConf
	Es                    struct {
		Address  []string
		Username string
		Password string
	}
	UserRPC zrpc.RpcClientConf
}
