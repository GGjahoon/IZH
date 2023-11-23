package xcode

import (
	"net/http"

	"github.com/GGjahoon/IZH/pkg/xcode/types"
)

// ErrHandler is a DIY error process http handler for recognize the XCode
// ErrHandler is used by applet-api http server as a router middware
// Hope that http response always is http.statusok and any include the real err
func ErrHandler(err error) (int, any) {
	code := CodeFromError(err)
	return http.StatusOK, types.Status{
		Code:    int32(code.Code()),
		Message: code.Message(),
	}
}
