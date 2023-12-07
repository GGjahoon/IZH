package main

import (
	"flag"
	"fmt"

	"github.com/GGjahoon/IZH/application/follow/rpc/internal/config"
	"github.com/GGjahoon/IZH/application/follow/rpc/internal/model"
	"github.com/GGjahoon/IZH/application/follow/rpc/internal/server"
	"github.com/GGjahoon/IZH/application/follow/rpc/internal/svc"
	"github.com/GGjahoon/IZH/application/follow/rpc/pb"
	"github.com/GGjahoon/IZH/pkg/orm"
	"github.com/GGjahoon/IZH/pkg/xcode/interceptors"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/follow.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	db := orm.MustNewMysql(&orm.Config{
		DSN:          c.DB.DataSource,
		MaxOpenConns: c.DB.MaxOpenConns,
		MaxIdleConns: c.DB.MaxIdleConns,
		MaxLifeTime:  c.DB.MaxLifeTime,
	})
	followModel := model.NewFollowModel(db.DB)
	followCountModel := model.NewFollowCountModel(db.DB)
	rds := redis.MustNewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})

	ctx := svc.NewServiceContext(c, db, followModel, followCountModel, rds)

	logx.DisableStat()
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterFollowServer(grpcServer, server.NewFollowServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(interceptors.ServerErrorInterceptor())
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
