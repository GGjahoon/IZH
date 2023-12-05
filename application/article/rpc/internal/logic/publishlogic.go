package logic

import (
	"context"
	"strconv"
	"time"

	"github.com/GGjahoon/IZH/application/article/rpc/internal/code"
	"github.com/GGjahoon/IZH/application/article/rpc/internal/model"
	"github.com/GGjahoon/IZH/application/article/rpc/internal/svc"
	"github.com/GGjahoon/IZH/application/article/rpc/internal/types"
	"github.com/GGjahoon/IZH/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishLogic) Publish(in *pb.PublishRequest) (*pb.PublishResponse, error) {
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if len(in.Title) == 0 {
		return nil, code.ArticleTitleCantEmpty
	}
	if len(in.Content) == 0 {
		return nil, code.ArticleContentCantEmpty
	}

	ret, err := l.svcCtx.ArticleModel.Insert(l.ctx, &model.Article{
		AuthorId:    in.UserId,
		Title:       in.Title,
		Content:     in.Content,
		Description: in.Description,
		Cover:       in.Cover,
		Status:      types.ArticleStatusVisible, //非正常逻辑，初次发布应为待审核
		PublishTime: time.Now(),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	})

	if err != nil {
		l.Logger.Error("Publish Insert req  : %v error : %v ", in, err)
		return nil, err
	}

	articleId, err := ret.LastInsertId()

	if err != nil {
		l.Logger.Error("Publish LastInsertID error :%v", err)
		return nil, err
	}
	//写入db后，代表该user的article列表发生了更新，需要同步更新redis中的数据，以避免读请求读到脏数据
	var (
		articleIdStr   = strconv.FormatInt(articleId, 10)
		likeNumKey     = articlesKey(in.UserId, types.SortLikeCount)
		publishTimeKey = articlesKey(in.UserId, types.SortPublishTime)
	)
	//在更新前先查询redis中是否缓存了该用户的文章列表
	exist, _ := l.svcCtx.BizRedis.Exists(likeNumKey)
	if exist {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, likeNumKey, 0, articleIdStr)
		if err != nil {
			logx.Errorf("Update Sort like_num ZaddCtx req:%v error:%v", in, err)
		}
	}

	exist, _ = l.svcCtx.BizRedis.Exists(publishTimeKey)
	if exist {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, publishTimeKey, time.Now().Unix(), articleIdStr)
		if err != nil {
			logx.Errorf("Update Sort publish_time ZaddCtx req:%v error:%v", in, err)
		}
	}

	return &pb.PublishResponse{ArticleId: articleId}, nil
}
