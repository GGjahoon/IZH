package main

import (
	"context"
	"flag"

	"github.com/GGjahoon/IZH/application/article/mq/internal/config"
	"github.com/GGjahoon/IZH/application/article/mq/internal/logic"
	"github.com/GGjahoon/IZH/application/article/mq/internal/model"
	"github.com/GGjahoon/IZH/application/article/mq/internal/svc"
	"github.com/GGjahoon/IZH/application/user/rpc/user"
	"github.com/GGjahoon/IZH/pkg/es"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

var configFile = flag.String("f", "etc/mq.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	err := c.ServiceConf.SetUp()
	if err != nil {
		panic(err)
	}
	logx.DisableStat()
	// begin to initialize servicecontext
	conn := sqlx.NewMysql(c.DataSource)
	articleModel := model.NewArticleModel(conn)

	rds, err := redis.NewRedis(redis.RedisConf{
		Host: c.BizRedisConf.Host,
		Pass: c.BizRedisConf.Pass,
		Type: c.BizRedisConf.Type,
	})
	if err != nil {
		panic(err)
	}
	userRPCClient := zrpc.MustNewClient(c.UserRPC)
	userRPC := user.NewUser(userRPCClient)

	es := es.MustNewEs(&es.Config{
		Address:  c.Es.Address,
		Username: c.Es.Username,
		Password: c.Es.Password,
	})

	svc := svc.NewServiceContext(c, articleModel, rds, userRPC, es)

	ctx := context.Background()

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	services := consumers(ctx, svc)

	for _, service := range services {
		serviceGroup.Add(service)
	}

	serviceGroup.Start()
}
func consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, logic.NewLikeNumLogic(ctx, svcCtx)),
		kq.MustNewQueue(svcCtx.Config.ArticleKqConsumerConf, logic.NewArticleLogic(ctx, svcCtx)),
	}
}
