package xcode

import (
	"fmt"
	"strconv"

	"github.com/GGjahoon/IZH/pkg/xcode/types"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var _ XCode = (*Status)(nil)

type Status struct {
	sts *types.Status
}

func Error(code Code) *Status {
	return &Status{
		sts: &types.Status{
			Code:    int32(code.Code()),
			Message: code.Message(),
		},
	}
}
func Errof(code Code, format string, args ...interface{}) *Status {
	code.msg = fmt.Sprintf(format, args...)
	return Error(code)
}
func (s *Status) Error() string {
	return s.Message()
}
func (s *Status) Code() int {
	return int(s.sts.Code)
}

func (s *Status) Message() string {
	if s.sts.Message == "" {
		return strconv.Itoa(int(s.sts.Code))
	}
	return s.sts.Message
}

func (s *Status) Details() []interface{} {
	if s == nil || s.sts == nil {
		return nil
	}
	details := make([]interface{}, 0, len(s.sts.Details))
	for _, d := range s.sts.Details {
		detail := &ptypes.DynamicAny{}
		if err := d.UnmarshalTo(detail); err != nil {
			details = append(details, err)
			continue
		}
		details = append(details, detail.Message)
	}
	return details
}

func (s *Status) WithDetails(msgs ...proto.Message) (*Status, error) {
	for _, msg := range msgs {
		anyMsg, err := anypb.New(msg)
		if err != nil {
			return s, err
		}
		s.sts.Details = append(s.sts.Details, anyMsg)
	}
	return s, nil
}
func (s *Status) Proto() *types.Status {
	return s.sts
}
func FromCode(code Code) *Status {
	return &Status{
		sts: &types.Status{
			Code:    int32(code.Code()),
			Message: code.Message(),
		},
	}
}
func FromProto(pbMsg proto.Message) XCode {
	msg, ok := pbMsg.(*types.Status)
	if ok {
		if len(msg.Message) == 0 || msg.Message == strconv.FormatInt(int64(msg.Code), 10) {
			return Code{code: int(msg.Code)}
		}
	}
	return Errof(ServerErr, "invalid proto message get %v", pbMsg)
}

// func toXCode(grpcStatus *status.Status) Code {
// 	grpcCode := grpcStatus.Code()
// 	switch grpcCode {
// 	case codes.OK:
// 		return OK
// 	case codes.InvalidArgument:
// 		return RequestErr
// 	case codes.NotFound:
// 		return NotFound
// 	case codes.PermissionDenied:
// 		return AccessDenied
// 	case codes.Unauthenticated:
// 		return Unauthorized
// 	case codes.ResourceExhausted:
// 		return LimitExceed
// 	case codes.DeadlineExceeded:
// 		return Deadline
// 	case codes.Unavailable:
// 		return ServiceUnavailable
// 	case codes.Unknown:
// 		return String(grpcStatus.Message())
// 	}
// 	return ServerErr
// }
