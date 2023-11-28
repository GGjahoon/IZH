package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/GGjahoon/IZH/application/article/rpc/internal/model"
	"github.com/GGjahoon/IZH/application/article/rpc/internal/svc"
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
	fmt.Println("start to publish")
	ret, err := l.svcCtx.ArticleModel.Insert(l.ctx, &model.Article{
		AuthorId:    in.UserId,
		Title:       in.Title,
		Content:     in.Content,
		Description: in.Description,
		Cover:       in.Cover,
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
	return &pb.PublishResponse{ArticleId: articleId}, nil
}
