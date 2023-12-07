package code

import "github.com/GGjahoon/IZH/pkg/xcode"

var (
	FollowUserIdEmpty   = xcode.NewCode(40001, "关注用户ID为空")
	FollowedUserIdEmpty = xcode.NewCode(40002, "被关注用户ID为空")
	CannotFollowSelf    = xcode.NewCode(40003, "不能关注自己")
	UserIdEmpty         = xcode.NewCode(40004, "用户id为空")
)
