package logic

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/GGjahoon/IZH/application/applet/internal/code"
	"github.com/GGjahoon/IZH/application/applet/internal/svc"
	"github.com/GGjahoon/IZH/application/applet/internal/types"
	"github.com/GGjahoon/IZH/application/user/rpc/user"
	"github.com/GGjahoon/IZH/pkg/encrypt"
	"github.com/GGjahoon/IZH/pkg/jwt"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	prefixActivation = "biz#activation#%s"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	// remove space and check request is correct or not
	req, err = validRequest(req)

	if err != nil {
		return nil, err
	}
	// check the verificationCode is correct or not
	err = checkVerificationCode(l.svcCtx.BizRedis, req.Mobile, req.VerificationCode)
	if err != nil {
		return nil, code.VerificationCodeMismatch
	}
	// encrypt the mobile
	mobile, err := encrypt.EncMobile(req.Mobile)
	if err != nil {
		logx.Errorf("EncMobile : %s error : %v", req.Mobile, err)
		return nil, err
	}
	// check this phone number is registered or not
	userRet, err := l.svcCtx.UserRpc.FindByMobile(l.ctx, &user.FindByMobileRequest{Mobile: mobile})
	if err != nil {
		return nil, err
	}
	if userRet != nil && userRet.UserId > 0 {
		return nil, code.MobileHasRegistered
	}
	//call user rpc service to save this user in db
	fmt.Println("begin to register user in db")
	regRet, err := l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterRequest{
		Username: req.Name,
		Mobile:   mobile,
	})
	fmt.Printf("%T %v", regRet.UserId, regRet.UserId)
	if err != nil {
		return nil, err
	}
	//if save user successed , create a token and append into the response
	jwtToken, err := jwt.BuildTokens(jwt.TokenOption{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]interface{}{
			"userId": regRet.UserId,
		},
	})
	if err != nil {
		logx.Errorf("build token err :%v", err)
		return nil, err
	}
	//if register successed and token build successed then del the verificationcode in cache
	err = delVerificationCache(req.Mobile, l.svcCtx.BizRedis)
	if err != nil {
		return nil, err
	}
	return &types.RegisterResponse{
		UserId: regRet.UserId,
		Token: types.Token{
			AccessToken:  jwtToken.AccessToken,
			AccessExpire: jwtToken.AccessExpire,
		},
	}, nil
}
func checkVerificationCode(rds *redis.Redis, mobile string, code string) error {
	cacheCode, err := getActivationCache(mobile, rds)
	if err != nil {
		return err
	}
	if cacheCode == "" {
		return errors.New("verification code expired")
	}
	if cacheCode != code {
		return errors.New("verification code is not correct")
	}
	return nil
}
func validRequest(req *types.RegisterRequest) (*types.RegisterRequest, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Mobile = strings.TrimSpace(req.Mobile)
	if len(req.Mobile) == 0 {
		return nil, code.RegisterMobileEmpty
	}
	// req.Password = strings.TrimSpace(req.Password)
	// if len(req.Password) == 0 {
	// 	return nil, code.RegisterPasswdEmpty
	// } else {
	// 	req.Password = encrypt.EncPassword(req.Password)
	// }
	req.VerificationCode = strings.TrimSpace(req.VerificationCode)
	if len(req.VerificationCode) == 0 {
		return nil, code.VerificationCodeEmpty
	}
	return req, nil
}
func delVerificationCache(mobile string, rds *redis.Redis) error {
	key := fmt.Sprintf(prefixActivation, mobile)
	_, err := rds.Del(key)
	return err
}
