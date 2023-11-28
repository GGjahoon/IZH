package logic

import (
	"context"
	"strings"

	"github.com/GGjahoon/IZH/application/applet/internal/code"
	"github.com/GGjahoon/IZH/application/applet/internal/svc"
	"github.com/GGjahoon/IZH/application/applet/internal/types"
	"github.com/GGjahoon/IZH/application/user/rpc/service"
	"github.com/GGjahoon/IZH/pkg/encrypt"
	"github.com/GGjahoon/IZH/pkg/jwt"
	"github.com/GGjahoon/IZH/pkg/xcode"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	//valid the parameters of req
	req, err = validLoginRequest(req, l.svcCtx.BizRedis)
	if err != nil {
		logx.Errorf("valid req error :%v", err)
		return nil, err
	}
	mobile, err := encrypt.EncMobile(req.Mobile)
	if err != nil {
		logx.Errorf("enc mobile error :%v", err)
		return nil, err
	}
	//find this user in db
	RPCrsp, err := l.svcCtx.UserRpc.FindByMobile(l.ctx, &service.FindByMobileRequest{
		Mobile: mobile,
	})
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, code.MobileIsNotRegistered
		}
		logx.Errorf("find by mobile error :%v", err)
		return nil, xcode.ServerErr
	}
	if RPCrsp == nil || RPCrsp.UserId == 0 {
		return nil, xcode.AccessDenied
	}
	jwtToken, err := jwt.BuildTokens(jwt.TokenOption{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]interface{}{
			"userId": RPCrsp.UserId,
		},
	})
	if err != nil {
		logx.Errorf("build token error :%v", err)
		return nil, xcode.ServerErr
	}
	_ = delVerificationCache(req.Mobile, l.svcCtx.BizRedis)
	return &types.LoginResponse{
		UserId: RPCrsp.UserId,
		Token: types.Token{
			AccessToken:  jwtToken.AccessToken,
			AccessExpire: jwtToken.AccessExpire,
		},
	}, nil
}
func validLoginRequest(req *types.LoginRequest, rds *redis.Redis) (*types.LoginRequest, error) {
	req.Mobile = strings.TrimSpace(req.Mobile)
	if len(req.Mobile) == 0 {
		return nil, code.LoginMobileEmpty
	}
	req.VerificationCode = strings.TrimSpace(req.VerificationCode)
	if len(req.VerificationCode) == 0 {
		return nil, code.VerificationCodeEmpty
	}
	err := checkVerificationCode(rds, req.Mobile, req.VerificationCode)
	if err != nil {
		return nil, code.VerificationCodeMismatch
	}
	return req, nil
}
