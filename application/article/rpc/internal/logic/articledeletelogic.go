package logic

import (
	"context"
	"fmt"

	"github.com/GGjahoon/IZH/application/article/rpc/internal/code"
	"github.com/GGjahoon/IZH/application/article/rpc/internal/svc"
	"github.com/GGjahoon/IZH/application/article/rpc/internal/types"
	"github.com/GGjahoon/IZH/application/article/rpc/pb"
	"github.com/GGjahoon/IZH/pkg/xcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDeleteLogic {
	return &ArticleDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleDeleteLogic) ArticleDelete(in *pb.ArticleDeletRequest) (*pb.ArticleDeletResponse, error) {
	fmt.Println("in article delete")
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if in.ArticleId <= 0 {
		return nil, code.ArticleIdInvalid
	}

	article, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		l.Logger.Errorf("ArticleDelete FindOne id: %d err : %v", in.ArticleId, err)
		return nil, err
	}
	if article.AuthorId != in.UserId {
		return nil, xcode.AccessDenied
	}
	//修改article status 为 user delete
	err = l.svcCtx.ArticleModel.UpdateArticleStatus(l.ctx, in.ArticleId, types.ArticleStatusUserDelete)
	if err != nil {
		l.Logger.Errorf("UpdateArticleStatus req : %v err : %v ", in, err)
		return nil, err
	}

	//从redis中删除该文章数据,操作redis可能出现错误，可忽略错误
	_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, articlesKey(in.UserId, types.SortLikeCount), in.ArticleId)
	if err != nil {
		l.Logger.Errorf("ZermCtx req : %v sortType: %d error : %v", in, types.SortLikeCount, err)
	}
	_, err = l.svcCtx.BizRedis.ZremCtx(l.ctx, articlesKey(in.UserId, types.SortPublishTime), in.ArticleId)
	if err != nil {
		l.Logger.Errorf("ZermCtx req : %v sortType: %d error : %v", in, types.SortLikeCount, err)
	}
	return &pb.ArticleDeletResponse{}, nil
}
