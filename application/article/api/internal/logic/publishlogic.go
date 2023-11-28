package logic

import (
	"context"
	"encoding/json"
	"fmt"

	code "github.com/GGjahoon/IZH/application/article/api/internal/Code"
	"github.com/GGjahoon/IZH/application/article/api/internal/svc"
	"github.com/GGjahoon/IZH/application/article/api/internal/types"
	"github.com/GGjahoon/IZH/application/article/rpc/pb"
	"github.com/GGjahoon/IZH/pkg/xcode"

	"github.com/zeromicro/go-zero/core/logx"
)

var minContentLen = 80

type PublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishLogic) Publish(req *types.PublishRequest) (resp *types.PublishResponse, err error) {
	fmt.Println("start to valid parameters")
	// valid the parameters of req
	if len(req.Content) == 0 {
		return nil, code.ArticleTitleEmpty
	}
	if len(req.Content) < minContentLen {
		return nil, code.ArticleCotentTooFewWords
	}
	if len(req.Cover) == 0 {
		return nil, code.ArticleCoverEmpty
	}
	//get user id
	fmt.Println("start to get userid")
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	fmt.Printf("userid is %v", userId)
	if err != nil {
		logx.Errorf("get user id error : %v", err)
		return nil, xcode.Nologin
	}
	fmt.Println("start to publish")
	pret, err := l.svcCtx.ArticleRpc.Publish(l.ctx, &pb.PublishRequest{
		UserId:      userId,
		Title:       req.Title,
		Content:     req.Content,
		Description: req.Description,
		Cover:       req.Cover,
	})
	if err != nil {
		logx.Errorf("user:%d publish article error :%v ", userId, err)
		return nil, err
	}
	return &types.PublishResponse{
		AriticleId: pret.ArticleId,
	}, nil
}
