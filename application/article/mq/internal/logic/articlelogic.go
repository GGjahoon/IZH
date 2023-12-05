package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/GGjahoon/IZH/application/article/mq/internal/svc"
	"github.com/GGjahoon/IZH/application/article/mq/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLogic {
	return &ArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleLogic) Consume(_, val string) error {
	logx.Infof("article logic Consume msg : %s", val)
	var msg *types.CanalArticleMsg
	err := json.Unmarshal([]byte(val), &msg)
	if err != nil {
		logx.Errorf("Consume val :%s error:%v ", val, err)
		return err
	}
	return l.articleOperate(l.ctx, msg)
}
func (l *ArticleLogic) articleOperate(ctx context.Context, msg *types.CanalArticleMsg) error {
	if len(msg.Data) == 0 {
		return nil
	}
	//从msg.Data中解析出数据
	for _, data := range msg.Data {
		status, _ := strconv.Atoi(data.Status)
		likeNum, _ := strconv.Atoi(data.LikeNum)
		publishTime, _ := time.ParseInLocation("2006-01-02 15:04:05", data.PublishTime, time.Local)

		publishTimeKey := articlesKey(data.AuthoId, 0)
		likeNumKey := articlesKey(data.AuthoId, 1)
		//根据不同status对该用户的文章列表进行操作，ArticleStatusVisible：加入该用户的文章列表
		//ArticleStatusUserDelete：从该用户的文章列表中删除
		switch status {
		// 若状态为可见，则添加到文章列表中
		case types.ArticleStatusVisible:
			//操作前先查看redis内是否有该用户的文章列表
			exist, _ := l.svcCtx.BizRedis.ExistsCtx(l.ctx, publishTimeKey)
			if exist {
				// 若文章列表存在则更新

				_, err := l.svcCtx.BizRedis.ZaddCtx(l.ctx, publishTimeKey, publishTime.Unix(), data.ID)
				if err != nil {
					l.Logger.Errorf("ZaddCtx key:%s req :%v error:%v", publishTimeKey, data, err)
				}
			}
			exist, _ = l.svcCtx.BizRedis.ExistsCtx(l.ctx, likeNumKey)
			if exist {
				_, err := l.svcCtx.BizRedis.ZaddCtx(l.ctx, likeNumKey, int64(likeNum), data.ID)
				if err != nil {
					l.Logger.Errorf("ZaddCtx key:%s req:%v error:%v", likeNumKey, data, err)
				}
			}
		//若状态为用户删除，则从文章列表中删除
		case types.ArticleStatusUserDelete:
			exist, _ := l.svcCtx.BizRedis.ExistsCtx(l.ctx, publishTimeKey)
			if exist {
				_, err := l.svcCtx.BizRedis.ZremCtx(l.ctx, publishTimeKey, data.ID)
				if err != nil {
					l.Logger.Errorf("ZremCtx key:%s req :%v error:%v", publishTimeKey, data, err)
				}
			}
			exist, _ = l.svcCtx.BizRedis.ExistsCtx(l.ctx, likeNumKey)
			if exist {
				_, err := l.svcCtx.BizRedis.ZremCtx(l.ctx, likeNumKey, data.ID)
				if err != nil {
					l.Logger.Errorf("ZremCtx key:%s req :%v error:%v", likeNumKey, data, err)
				}
			}
		}

	}
	return nil
}
func articlesKey(userId string, sortType int32) string {
	return fmt.Sprintf("biz#articles#%s#%d", userId, sortType)
}
