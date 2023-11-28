package logic

import (
	"context"
	"errors"

	"github.com/GGjahoon/IZH/application/user/rpc/internal/code"
	"github.com/GGjahoon/IZH/application/user/rpc/internal/model"
	"github.com/GGjahoon/IZH/application/user/rpc/internal/svc"
	"github.com/GGjahoon/IZH/application/user/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *service.RegisterRequest) (*service.RegisterResponse, error) {
	if len(in.Username) == 0 {
		return nil, code.RegisterNameEmpty
	}
	ret, err := l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		Username: in.Username,
		Mobile:   in.Mobile,
		Avatar:   in.Avatar,
	})
	if err != nil {
		logx.Errorf("Register req : %v error: %v", in, err)
		return nil, err
	}
	userId, err := ret.LastInsertId()
	if err != nil {
		logx.Errorf("LatstInsertId error : %v ", err)
		return nil, errors.New("cannot get user id")
	}

	return &service.RegisterResponse{UserId: userId}, nil
}
