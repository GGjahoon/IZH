package main

import (
	"flag"
	"fmt"

	"github.com/GGjahoon/IZH/application/article/api/internal/config"
	"github.com/GGjahoon/IZH/application/article/api/internal/handler"
	"github.com/GGjahoon/IZH/application/article/api/internal/svc"
	"github.com/GGjahoon/IZH/application/article/rpc/article"
	"github.com/GGjahoon/IZH/application/user/rpc/user"
	"github.com/GGjahoon/IZH/pkg/xcode"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/zrpc"
)

const (
	defaultOssConnectTimeout   = 1
	defaultOssReadWriteTimeout = 3
)

var configFile = flag.String("f", "etc/applet-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	//Initialize OssClient here
	if c.Oss.ConnTimeout == 0 {
		c.Oss.ConnTimeout = defaultOssConnectTimeout
	}
	if c.Oss.ReadWriteTimeout == 0 {
		c.Oss.ReadWriteTimeout = defaultOssReadWriteTimeout
	}
	oc, err := oss.New(c.Oss.Endpoint, c.Oss.AccessKey, c.Oss.AccessKeySecret, oss.Timeout(c.Oss.ConnTimeout, c.Oss.ReadWriteTimeout))
	if err != nil {
		panic(err)
	}
	//Initialize UserRpcClient here
	userClient := zrpc.MustNewClient(c.UserRPC)
	userRpcClient := user.NewUser(userClient)
	//Initialize ArticleRpcClient here
	articleClient := zrpc.MustNewClient(c.ArticleRPC)
	articleRpcClient := article.NewArticle(articleClient)

	ctx := svc.NewServiceContext(c, oc, userRpcClient, articleRpcClient)
	handler.RegisterHandlers(server, ctx)
	httpx.SetErrorHandler(xcode.ErrHandler)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
