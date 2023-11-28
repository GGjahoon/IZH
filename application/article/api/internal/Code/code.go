package code

import "github.com/GGjahoon/IZH/pkg/xcode"

var (
	GetFileErr               = xcode.NewCode(30000, "文件获取失败，清添加封面图片")
	GetBucketErr             = xcode.NewCode(30001, "获取bucket实例失败")
	PutBucketErr             = xcode.NewCode(30002, "上传文件失败")
	GetObjDetailErr          = xcode.NewCode(30003, "获取对象详细信息失败")
	ArticleTitleEmpty        = xcode.NewCode(30004, "文章标题为空")
	ArticleCotentTooFewWords = xcode.NewCode(30005, "文章内容字数太少")
	ArticleCoverEmpty        = xcode.NewCode(30006, "未上传文章封面")
)
