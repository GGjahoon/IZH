package main

import (
	"context"
	"flag"

	"github.com/GGjahoon/IZH/application/like/mq/internal/config"
	"github.com/GGjahoon/IZH/application/like/mq/internal/logic"
	"github.com/GGjahoon/IZH/application/like/mq/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/like.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	//load etc/like.yaml into config
	conf.MustLoad(*configFile, &c)

	//generate a new service context
	svc := svc.NewServiceContext(c)
	ctx := context.Background()
	//generate a new service group
	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()
	services := logic.Consumes(ctx, svc)
	for _, mq := range services {
		serviceGroup.Add(mq)
	}

	serviceGroup.Start()
}
