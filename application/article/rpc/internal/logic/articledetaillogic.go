package logic

import (
	"context"

	"github.com/GGjahoon/IZH/application/article/rpc/internal/code"
	"github.com/GGjahoon/IZH/application/article/rpc/internal/svc"
	"github.com/GGjahoon/IZH/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ArticleDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDetailLogic {
	return &ArticleDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleDetailLogic) ArticleDetail(in *pb.ArticleDetailRequest) (*pb.ARticleDetailResponse, error) {
	if in.ArticleId <= 0 {
		return nil, code.ArticleIdInvalid
	}
	article, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return &pb.ARticleDetailResponse{}, nil
		}
		return nil, err
	}
	return &pb.ARticleDetailResponse{
		Article: &pb.ArticleItem{
			Id:    article.Id,
			Title: article.Title,
		},
	}, nil
}
