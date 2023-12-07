package logic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/GGjahoon/IZH/application/follow/rpc/internal/code"
	"github.com/GGjahoon/IZH/application/follow/rpc/internal/model"
	"github.com/GGjahoon/IZH/application/follow/rpc/internal/svc"
	"github.com/GGjahoon/IZH/application/follow/rpc/internal/types"
	"github.com/GGjahoon/IZH/application/follow/rpc/pb"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowLogic) Follow(in *pb.FollowRequest) (*pb.FollowResponse, error) {
	fmt.Println("in follow")
	if in.UserId == 0 {
		return nil, code.FollowUserIdEmpty
	}
	if in.FollowedUserId == 0 {
		return nil, code.FollowedUserIdEmpty
	}
	if in.UserId == in.FollowedUserId {
		return nil, code.CannotFollowSelf
	}

	follow, err := l.svcCtx.FollowModel.FindByUserIdAndFollowedUserID(l.ctx, in.UserId, in.FollowedUserId)

	if err != nil {
		l.Logger.Errorf("[Follow] FindByUserIdAndFollowedUserID err :%v req :%v", err, in)
		return nil, err
	}
	if follow != nil && follow.FollowStatus == types.FollowStatusFollow {
		return &pb.FollowResponse{}, nil
	}
	fmt.Println("begin transaction")
	//事务(需要将修改关注状态/创建关注记录和关注数++在一个事务下进行，避免两张表数据不一致)
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if follow != nil {
			err = model.NewFollowModel(tx).UpdateFields(l.ctx, follow.ID, map[string]interface{}{
				"follow_status": types.FollowStatusFollow,
			})
		} else {
			err = model.NewFollowModel(tx).Insert(l.ctx, &model.Follow{
				UserId:         in.UserId,
				FollowedUserId: in.FollowedUserId,
				FollowStatus:   types.FollowStatusFollow,
				CreateTime:     time.Now(),
				UpdateTime:     time.Now(),
			})
		}

		if err != nil {
			return err
		}
		//关注成功后将follow_count ++
		err = model.NewFollowCountModel(tx).IncrFollowCount(l.ctx, in.UserId)
		if err != nil {
			return err
		}
		return model.NewFollowCountModel(tx).IncrFansCount(l.ctx, in.FollowedUserId)
	})
	if err != nil {
		l.Logger.Errorf("[Follow] Transaction error :%v", err)
		return nil, err
	}
	//同步写入到缓存
	//check exist or not
	isFollowExist, err := l.svcCtx.BizRedis.ExistsCtx(l.ctx, userFollowKey(in.UserId))
	if err != nil {
		l.Logger.Errorf("[Follow] Reids Exists error:%v", err)
		return nil, err
	}
	if isFollowExist {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, userFollowKey(in.UserId), time.Now().Unix(), strconv.FormatInt(in.FollowedUserId, 10))
		if err != nil {
			l.Logger.Errorf("[Follow] Redis Zadd error: %v ", err)
			return nil, err
		}
		//删除一部分数据，即redis内的用户关注列表不能存储所有的关注用户id
		_, err = l.svcCtx.BizRedis.ZremrangebyrankCtx(l.ctx, userFollowKey(in.UserId), 0, -(types.CacheMaxFollowCount + 1))
		if err != nil {
			l.Logger.Errorf("[Follow] Reids ZremrangebyrankCtx err : %v ", err)
		}
	}

	isFansExist, err := l.svcCtx.BizRedis.ExistsCtx(l.ctx, userFansKey(in.FollowedUserId))
	if err != nil {
		l.Logger.Errorf("[Follow] Reids Exists error:%v", err)
		return nil, err
	}
	if isFansExist {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, userFansKey(in.FollowedUserId), time.Now().Unix(), strconv.FormatInt(in.UserId, 10))
		if err != nil {
			l.Logger.Errorf("[Follow] Redis Zadd error: %v ", err)
			return nil, err
		}
		_, err = l.svcCtx.BizRedis.ZremrangebyrankCtx(l.ctx, userFansKey(in.FollowedUserId), 0, -(types.CacheMaxFansCount + 1))
		if err != nil {
			l.Logger.Errorf("[Follow] Reids ZremrangebyrankCtx err : %v ", err)
		}
	}
	return &pb.FollowResponse{}, nil
}

func userFollowKey(userId int64) string {
	return fmt.Sprintf("biz#user#follow#%d", userId)
}
func userFansKey(userID int64) string {
	return fmt.Sprintf("biz#user#fans#%d", userID)
}
