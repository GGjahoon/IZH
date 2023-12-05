package logic

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/GGjahoon/IZH/application/article/mq/internal/svc"
	"github.com/GGjahoon/IZH/application/article/mq/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleLikeNumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeNumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLikeNumLogic {
	return &ArticleLikeNumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *ArticleLikeNumLogic) Consume(_, v string) error {
	logx.Infof("article like num logic Consume msg : %s", v)
	var msg *types.CanalLikeMsg
	err := json.Unmarshal([]byte(v), &msg)
	if err != nil {
		l.Logger.Errorf("Consume val: %s error : %v", v, err)
		return err
	}
	return l.updateLikeNum(l.ctx, msg)
}

func (l *ArticleLikeNumLogic) updateLikeNum(ctx context.Context, msg *types.CanalLikeMsg) error {
	if len(msg.Data) == 0 {
		return nil
	}
	for _, d := range msg.Data {
		if d.BizID != types.ArticleBizID {
			continue
		}
		id, err := strconv.ParseInt(d.ObjID, 10, 64)
		if err != nil {
			l.Logger.Errorf("ParseInt %s error %v", d.ObjID, err)
			continue
		}
		likeNum, err := strconv.ParseInt(d.LikeNum, 10, 64)
		if err != nil {
			l.Logger.Errorf("ParseInt %s err %v", d.LikeNum, err)
			continue
		}
		err = l.svcCtx.ArticleModel.UpdateLikeNum(ctx, id, likeNum)
		if err != nil {
			l.Logger.Errorf("UpdateLikeNum id :%d like:%d", id, likeNum)
		}
	}
	return nil
}
