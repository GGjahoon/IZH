package main

import (
	"flag"
	"fmt"

	"github.com/GGjahoon/IZH/application/article/rpc/internal/config"
	"github.com/GGjahoon/IZH/application/article/rpc/internal/model"
	"github.com/GGjahoon/IZH/application/article/rpc/internal/server"
	"github.com/GGjahoon/IZH/application/article/rpc/internal/svc"
	"github.com/GGjahoon/IZH/application/article/rpc/pb"
	"github.com/GGjahoon/IZH/pkg/xcode/interceptors"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/article.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	//initialize article model
	conn := sqlx.NewMysql(c.DataSource)
	articleModel := model.NewArticleModel(conn, c.CacheRedis)

	//initialize redis
	rds, err := redis.NewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})
	if err != nil {
		panic(err)
	}

	//initialize service context
	ctx := svc.NewServiceContext(c, articleModel, rds)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterArticleServer(grpcServer, server.NewArticleServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	// add a custom interceptor
	s.AddUnaryInterceptors(interceptors.ServerErrorInterceptor())
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
