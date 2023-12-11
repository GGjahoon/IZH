package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/GGjahoon/IZH/application/article/mq/internal/svc"
	"github.com/GGjahoon/IZH/application/article/mq/internal/types"
	"github.com/GGjahoon/IZH/application/user/rpc/user"
	"github.com/elastic/go-elasticsearch/v8/esutil"
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
	return l.articleOperate(msg)
}
func (l *ArticleLogic) articleOperate(msg *types.CanalArticleMsg) error {
	if len(msg.Data) == 0 {
		return nil
	}
	var esData []*types.ArticleEsMsg
	//从msg.Data中解析出数据
	for _, data := range msg.Data {
		status, _ := strconv.Atoi(data.Status)
		likeNum, _ := strconv.ParseInt(data.LikeNum, 10, 64)
		articleId, _ := strconv.ParseInt(data.ID, 10, 64)
		authorId, _ := strconv.ParseInt(data.AuthorId, 10, 64)

		publishTime, _ := time.ParseInLocation("2006-01-02 15:04:05", data.PublishTime, time.Local)
		publishTimeKey := articlesKey(data.AuthorId, 0)
		likeNumKey := articlesKey(data.AuthorId, 1)

		//redis：
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

		//ES：
		// 找到用户，聚合数据，写入ES
		user, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{
			UserId: authorId,
		})
		if err != nil {
			l.Logger.Errorf("Find By Id : %d error : %v", authorId, err)
			return err
		}
		esData = append(esData, &types.ArticleEsMsg{
			ArticleID:   articleId,
			Title:       data.Title,
			Content:     data.Content,
			Description: data.Description,
			AuthorId:    authorId,
			AuthorName:  user.Username,
			Status:      status,
			LikeNum:     likeNum,
		})
	}

	//批量写入ES
	err := l.BatchUpSertToEs(l.ctx, esData)
	if err != nil {
		l.Logger.Errorf("l.BatchUpSertToEs data : %v error : %v ", esData, err)
	}

	return nil
}
func articlesKey(userId string, sortType int32) string {
	return fmt.Sprintf("biz#articles#%s#%d", userId, sortType)
}

func (l *ArticleLogic) BatchUpSertToEs(ctx context.Context, data []*types.ArticleEsMsg) error {
	if len(data) == 0 {
		return nil
	}
	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: l.svcCtx.Es.Client,
		Index:  "article-index",
	})
	if err != nil {
		return err
	}
	fmt.Println(ctx)
	for _, d := range data {
		v, err := json.Marshal(d)
		if err != nil {
			return err
		}
		payload := fmt.Sprintf(`{"doc":%s,"doc_as_upsert":true}`, string(v))
		fmt.Println("begin add ")
		err = bulkIndexer.Add(ctx, esutil.BulkIndexerItem{
			Action:     "update",
			DocumentID: fmt.Sprintf("%d", d.ArticleID),
			Body:       strings.NewReader(payload),
			OnSuccess: func(ctx context.Context, bii esutil.BulkIndexerItem, biri esutil.BulkIndexerResponseItem) {
				fmt.Println("写入成功")
			},
			OnFailure: func(ctx context.Context, bii esutil.BulkIndexerItem, biri esutil.BulkIndexerResponseItem, err error) {
				fmt.Println("add error : ", err)
			},
		})
		fmt.Println("finsh add ")
		if err != nil {
			return err
		}
	}
	return bulkIndexer.Close(ctx)
}
