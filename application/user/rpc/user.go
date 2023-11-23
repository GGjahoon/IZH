package main

import (
	"flag"
	"fmt"

	"github.com/GGjahoon/IZH/application/user/rpc/internal/config"
	"github.com/GGjahoon/IZH/application/user/rpc/internal/model"
	"github.com/GGjahoon/IZH/application/user/rpc/internal/server"
	"github.com/GGjahoon/IZH/application/user/rpc/internal/svc"
	"github.com/GGjahoon/IZH/application/user/rpc/service"
	"github.com/GGjahoon/IZH/pkg/xcode/interceptors"

	"github.com/zeromicro/go-zero/core/conf"
	cs "github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	conn := sqlx.NewMysql(c.DataSource)
	userModel := model.NewUserModel(conn, c.CacheRedis)
	ctx := svc.NewServiceContext(c, userModel)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		service.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == cs.DevMode || c.Mode == cs.TestMode {
			reflection.Register(grpcServer)
		}
	})
	//user the diy error interceptor
	s.AddUnaryInterceptors(interceptors.ServerErrorInterceptor())
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
