package logic

import (
	"context"
	"fmt"

	"github.com/GGjahoon/IZH/application/user/rpc/internal/svc"
	"github.com/GGjahoon/IZH/application/user/rpc/service"
	"github.com/GGjahoon/IZH/pkg/xcode"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type FindByMobileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByMobileLogic {
	return &FindByMobileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByMobileLogic) FindByMobile(in *service.FindByMobileRequest) (*service.FindByMobileResponse, error) {
	fmt.Println("start to find by mobile")
	user, err := l.svcCtx.UserModel.FindByMobile(l.ctx, in.Mobile)
	if err != nil {
		logx.Errorf("FindByMobile : %s error : %v", in.Mobile, err)
		if err == sqlx.ErrNotFound {
			return nil, xcode.NotFound
		}
		return nil, xcode.FindByMobileErr
	}

	return &service.FindByMobileResponse{
		UserId:   user.Id,
		Username: user.Username,
		Avatar:   user.Avatar,
	}, nil
}
