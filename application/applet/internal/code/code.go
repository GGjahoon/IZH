package code

import "github.com/GGjahoon/IZH/pkg/xcode"

var (
	RegisterMobileEmpty      = xcode.NewCode(10001, "注册手机号码不能为空")
	VerificationCodeEmpty    = xcode.NewCode(100002, "验证码不能为空")
	MobileHasRegistered      = xcode.NewCode(100003, "手机号已被注册")
	LoginMobileEmpty         = xcode.NewCode(100003, "手机号不能为空")
	RegisterPasswdEmpty      = xcode.NewCode(100004, "密码不能为空")
	VerificationCodeMismatch = xcode.NewCode(100005, "验证码错误")
)
