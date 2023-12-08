package logic

import (
	"context"
	"strconv"

	"github.com/GGjahoon/IZH/application/follow/rpc/internal/code"
	"github.com/GGjahoon/IZH/application/follow/rpc/internal/model"
	"github.com/GGjahoon/IZH/application/follow/rpc/internal/svc"
	"github.com/GGjahoon/IZH/application/follow/rpc/internal/types"
	"github.com/GGjahoon/IZH/application/follow/rpc/pb"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnFollowLogic {
	return &UnFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnFollowLogic) UnFollow(in *pb.UnFollowRequest) (*pb.UnFollowResponse, error) {
	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}
	if in.FollowedUserId == 0 {
		return nil, code.FollowUserIdEmpty
	}
	follow, err := l.svcCtx.FollowModel.FindByUserIdAndFollowedUserID(l.ctx, in.UserId, in.FollowedUserId)
	if err != nil {
		l.Logger.Errorf("[UnFollow] FollowModel.FindByUserIdAndFollowedUserID err : %v req : %v", err, in)
		return nil, err
	}
	if follow == nil {
		return &pb.UnFollowResponse{}, nil
	}
	if follow.FollowStatus == types.FollowStatusUnfollow {
		return &pb.UnFollowResponse{}, nil
	}

	//开启事务
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		err := model.NewFollowModel(tx).UpdateFields(l.ctx, follow.ID, map[string]interface{}{
			"follow_status": types.FollowStatusUnfollow,
		})
		if err != nil {
			return err
		}
		err = model.NewFollowCountModel(tx).DecrFollowCount(l.ctx, in.UserId)
		if err != nil {
			return err
		}
		return model.NewFollowCountModel(tx).DecrFansCount(l.ctx, in.FollowedUserId)
	})
	if err != nil {
		l.Logger.Errorf("[UnFollow] Transaction error : %v", err)
		return nil, err
	}

	//同步更新redis，从关注列表和粉丝列表中删除对应数据
	_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, userFollowKey(in.UserId), strconv.FormatInt(in.FollowedUserId, 10))
	if err != nil {
		l.Logger.Errorf("[UnFollow] BizRedis ZremCtx Follower list error : %v", err)
		return nil, err
	}

	_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, userFansKey(in.FollowedUserId), strconv.FormatInt(in.UserId, 10))
	if err != nil {
		l.Logger.Errorf("[UnFollow] BizRedis ZremCtx fans list error : %v", err)
		return nil, err
	}
	return &pb.UnFollowResponse{}, nil
}
