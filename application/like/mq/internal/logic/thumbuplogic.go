package logic

import (
	"context"
	"fmt"

	"github.com/GGjahoon/IZH/application/like/mq/internal/svc"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

// ThumbupLogic contains ctx and service context
type ThumbupLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
	logx.Logger
}

func NewThumbupLogic(ctx context.Context, svc *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		ctx:    ctx,
		svc:    svc,
		Logger: logx.WithContext(ctx),
	}
}

// ThumbupLogic's Consume method
func (l *ThumbupLogic) Consume(key, value string) error {
	fmt.Printf("get key : %s value :%s\n", key, value)
	return nil
}

// 将服务注册到service.Service interface中
// kq.MustNewQueue入参:handler kq.ConsumeHandler(包含Consume method 的interface)
func Consumes(ctx context.Context, svc *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svc.Config.KqConsumerConf, NewThumbupLogic(ctx, svc)),
	}
}
