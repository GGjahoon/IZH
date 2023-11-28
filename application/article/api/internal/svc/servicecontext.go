package svc

import (
	"github.com/GGjahoon/IZH/application/article/api/internal/config"
	"github.com/GGjahoon/IZH/application/article/rpc/article"
	"github.com/GGjahoon/IZH/application/user/rpc/user"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type ServiceContext struct {
	Config     config.Config
	OssClient  *oss.Client
	UserRpc    user.User
	ArticleRpc article.Article
}

func NewServiceContext(c config.Config,
	ossClient *oss.Client,
	userRpc user.User,
	articleRpc article.Article,
) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		OssClient:  ossClient,
		UserRpc:    userRpc,
		ArticleRpc: articleRpc,
	}
}
