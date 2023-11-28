package logic

import (
	"context"
	"encoding/json"

	"github.com/GGjahoon/IZH/application/applet/internal/svc"
	"github.com/GGjahoon/IZH/application/applet/internal/types"
	"github.com/GGjahoon/IZH/application/user/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (resp *types.UserInfoResponse, err error) {
	userId, err := l.ctx.Value(types.UserIDKey).(json.Number).Int64()
	if err != nil {
		return nil, err
	}
	if userId == 0 {
		return nil, nil
	}
	user, err := l.svcCtx.UserRpc.FindById(l.ctx, &service.FindByIdRequest{
		UserId: userId,
	})
	if err != nil {
		logx.Errorf("FindById userID: %d error : %v", userId, err)
		return nil, err
	}

	return &types.UserInfoResponse{
		UserId:   user.UserId,
		Username: user.Username,
		Avatar:   user.Avatar,
	}, nil
}
