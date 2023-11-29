package logic

import (
	"context"
	"encoding/json"

	"github.com/GGjahoon/IZH/application/like/rpc/internal/svc"
	"github.com/GGjahoon/IZH/application/like/rpc/internal/types"
	"github.com/GGjahoon/IZH/application/like/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
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
	// todo: add the real Thumbup logic
	// one: check this user has thumbup already
	// two : caculate the times of thumbup and thumbdown count of  current content

	// generate a kafka msg use to send into kafka queue
	msg := types.ThumbupMsg{
		BizId:    in.BizId,
		ObjId:    in.ObjId,
		UserId:   in.UserId,
		LikeType: in.LikeType,
	}
	// send the kafka message asynchronously
	threading.GoSafe(func() {
		data, err := json.Marshal(msg)
		if err != nil {
			l.Logger.Errorf("[Thumbup]  marshal msg : %v error: %v", msg, err)
			return
		}
		err = l.svcCtx.KqPusherClient.Push(string(data))
		if err != nil {
			l.Logger.Errorf("[Thumbup] kq push data : %s error : %v", string(data), err)
			return
		}
	})
	return &service.ThumbupResponse{
		LikeNum: 6,
	}, nil
}
