package logic

import (
	"context"
	"fmt"
	"net/http"
	"time"

	code "github.com/GGjahoon/IZH/application/article/api/internal/Code"
	"github.com/GGjahoon/IZH/application/article/api/internal/svc"
	"github.com/GGjahoon/IZH/application/article/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const maxFileSize = 10 << 20 //10MB

type UploadCoverLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadCoverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadCoverLogic {
	return &UploadCoverLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadCoverLogic) UploadCover(req *http.Request) (resp *types.UploadCoverResponse, err error) {
	// parse the http request body
	req.ParseMultipartForm(maxFileSize)
	// get the picture file which the key is "cover" in http request
	file, handler, err := req.FormFile("cover")
	if err != nil {
		logx.Errorf("get file failed :%v", err)
		return nil, code.GetFileErr
	}
	defer file.Close()
	bucket, err := l.svcCtx.OssClient.Bucket(l.svcCtx.Config.Oss.BucketName)
	if err != nil {
		logx.Errorf("get bucket failed :%v", err)
		return nil, code.GetBucketErr
	}
	objectKey := genFilename(handler.Filename)
	err = bucket.PutObject(objectKey, file)
	if err != nil {
		logx.Errorf("push file failed :%v", err)
		return nil, code.PutBucketErr
	}
	return &types.UploadCoverResponse{CoverUrl: genFileUrl(objectKey)}, nil

}
func genFilename(filename string) string {
	return fmt.Sprintf("%d_%s", time.Now().UnixMilli(), filename)
}
func genFileUrl(objectKey string) string {
	return fmt.Sprintf("https://article-cover.oss-cn-beijing.aliyuncs.com/%s", objectKey)
}
