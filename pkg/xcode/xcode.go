package xcode

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
)

// XCode is a DIY Error interface
type XCode interface {
	Error() string
	Code() int
	Message() string
	Details() []interface{}
}

// Code is athe real implementation of XCode
type Code struct {
	code int
	msg  string
}

// NewCode returns a new Code include all function of XCode
func NewCode(code int, msg string) XCode {
	return &Code{
		code: code,
		msg:  msg,
	}
}
func (code *Code) Error() string {
	return code.Message()
}
func (code *Code) Code() int {
	return code.code
}
func (code *Code) Message() string {
	return code.msg
}
func (code *Code) Details() []interface{} {
	return nil
}

// CodeFromError convert the error into XCode interface to response to user client
func CodeFromError(err error) XCode {
	err = errors.Cause(err)
	//XCode与err均有Error()方法，将err断言为XCode
	if code, ok := err.(XCode); ok {
		return code
	}
	switch err {
	case context.Canceled:
		return Canceled
	case context.DeadlineExceeded:
		return Deadline
	}
	return ServerErr
}

func String(s string) XCode {
	if len(s) == 0 {
		return OK
	}
	code, err := strconv.Atoi(s)
	if err != nil {
		return ServerErr
	}
	return &Code{code: code}

}
