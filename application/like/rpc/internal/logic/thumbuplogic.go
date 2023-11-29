package logic

import (
	"context"

	"github.com/GGjahoon/IZH/application/like/rpc/internal/svc"
	"github.com/GGjahoon/IZH/application/like/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type ThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ThumbupLogic) Thumbup(in *service.ThumbupRequest) (*service.ThumbupResponse, error) {
	// todo: add your logic here and delete this line

	return &service.ThumbupResponse{}, nil
}
