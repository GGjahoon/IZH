package logic

import (
	"context"

	"github.com/GGjahoon/IZH/application/user/rpc/internal/svc"
	"github.com/GGjahoon/IZH/application/user/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type SmsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSmsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SmsLogic {
	return &SmsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SmsLogic) Sms(in *service.SendSmsRequest) (*service.SendSmsResponse, error) {
	// todo: add your logic here and delete this line

	return &service.SendSmsResponse{}, nil
}
