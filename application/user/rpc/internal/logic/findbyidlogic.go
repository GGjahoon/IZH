package logic

import (
	"context"
	"errors"

	"github.com/GGjahoon/IZH/application/user/rpc/internal/svc"
	"github.com/GGjahoon/IZH/application/user/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type FindByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByIdLogic {
	return &FindByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByIdLogic) FindById(in *service.FindByIdRequest) (*service.FindByIdResponse, error) {
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		logx.Errorf("FindBy userId : %d err : %v", in.UserId, err)
		if err == sqlx.ErrNotFound {
			return nil, errors.New("the user is not found in db")
		}
		return nil, errors.New("db internal error")
	}

	return &service.FindByIdResponse{
		UserId:   user.Id,
		Username: user.Username,
		Mobile:   user.Mobile,
		Avatar:   user.Avatar,
	}, nil
}
