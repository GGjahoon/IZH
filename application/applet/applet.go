package main

import (
	"flag"
	"fmt"

	"github.com/GGjahoon/IZH/application/applet/internal/config"
	"github.com/GGjahoon/IZH/application/applet/internal/handler"
	"github.com/GGjahoon/IZH/application/applet/internal/svc"
	"github.com/GGjahoon/IZH/pkg/xcode"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/applet-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	//设置自定义的ErrorHandler
	httpx.SetErrorHandler(xcode.ErrHandler)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
