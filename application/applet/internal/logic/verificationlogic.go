package logic

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/GGjahoon/IZH/application/applet/internal/svc"
	"github.com/GGjahoon/IZH/application/applet/internal/types"
	"github.com/GGjahoon/IZH/application/user/rpc/service"
	"github.com/GGjahoon/IZH/pkg/util"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	prefixVerificationCount = "biz#activation#count#%s"
	verificationLimitPerDay = 10
	expireActivation        = 60 * 30
)

type VerificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerificationLogic {
	return &VerificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerificationLogic) Verification(req *types.VerificationRequset) (resp *types.VerificationResponse, err error) {
	// get the count of this mobile's getting verification code
	count, err := l.getVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("getVerificationCount mobile : %s error : %v", req.Mobile, err)
	}
	if count > verificationLimitPerDay {
		return nil, errors.New("pass the max get verification code time")
	}
	// if the mobile get code before,give the code in cache,code does not change in 30 minutes
	code, err := getActivationCache(req.Mobile, l.svcCtx.BizRedis)
	if err != nil {
		logx.Errorf("getActivationCache mobile :%s error : %v", req.Mobile, err)
	}
	if len(code) == 0 {
		code = util.RandomNumber(6)
	}
	_, err = l.svcCtx.UserRpc.Sms(l.ctx, &service.SendSmsRequest{
		Mobile: req.Mobile,
	})
	if err != nil {
		logx.Errorf("send sms mobile: %s error:%v", req.Mobile, err)
		return nil, err
	}
	err = saveActivationCache(req.Mobile, code, l.svcCtx.BizRedis)
	if err != nil {
		logx.Errorf("saveActivationCache Mobile:%s error :%v", req.Mobile, err)
		return nil, err
	}
	err = l.incrVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("incrVerificationCount Mobile:%s error :%v", req.Mobile, err)
		return nil, err
	}

	return &types.VerificationResponse{}, nil
}

// incrVerificationCount increase the count of get verification code for one mobile and sset the expired time
func (l *VerificationLogic) incrVerificationCount(mobile string) error {
	//key := prefixVerificationCount + mobile
	key := fmt.Sprintf(prefixVerificationCount, mobile)
	_, err := l.svcCtx.BizRedis.Incr(key)
	if err != nil {
		return err
	}
	return l.svcCtx.BizRedis.Expireat(key, util.EndOfDay(time.Now()).Unix())
}

// getVerificationCount return the count of one mobile get the verification code
func (l *VerificationLogic) getVerificationCount(mobile string) (int, error) {
	//key := prefixActivation + mobile
	key := fmt.Sprintf(prefixVerificationCount, mobile)
	val, err := l.svcCtx.BizRedis.Get(key)
	if err != nil {
		return 0, err
	}
	if len(val) == 0 {
		return 0, nil
	}
	return strconv.Atoi(val)
}

// getActivationCache return the verification code in redis
func getActivationCache(mobile string, rds *redis.Redis) (string, error) {
	//key := prefixActivation + mobile
	key := fmt.Sprintf(prefixActivation, mobile)
	return rds.Get(key)
}

// saveActivationCache save the verification code in redis cache
func saveActivationCache(mobile string, code string, rds *redis.Redis) error {
	//key := prefixActivation + mobile
	key := fmt.Sprintf(prefixActivation, mobile)
	return rds.Setex(key, code, expireActivation)
}

// func delActivationCache(mobile string, rds *redis.Redis) error {
// 	key := fmt.Sprintf(prefixActivation, mobile)
// 	_, err := rds.Del(key)
// 	return err
// }
