package code

import "github.com/GGjahoon/IZH/pkg/xcode"

var (
	SortTypeInvalid         = xcode.NewCode(60001, "排序类型无效")
	UserIdInvalid           = xcode.NewCode(60002, "用户ID无效")
	ArticleTitleCantEmpty   = xcode.NewCode(60003, "请输入文章标题")
	ArticleContentCantEmpty = xcode.NewCode(60004, "请输入文章内容")
	ArticleIdInvalid        = xcode.NewCode(60005, "文章ID无效")
)
