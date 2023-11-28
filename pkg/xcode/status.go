package xcode

import (
	"context"
	"fmt"
	"strconv"

	"github.com/GGjahoon/IZH/pkg/xcode/types"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Status struct {
	sts *types.Status
}

// NewStatus return a new Status generate by Code and it has all method of XCode interface
func NewStatus(code Code) XCode {
	return &Status{
		sts: &types.Status{
			Code:    int32(code.code),
			Message: code.Message(),
		},
	}
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
			details = append(details, detail)
			continue
		}
		details = append(details, detail.Message)
	}
	return details
}

// Here is try to convert the Code error into grpc status in server

// in this grpc logic, some return error is XCode format , before return to grpc Client,
// must convert the err (which is XCode format) into gRPC status
func FromErr(err error) *status.Status {
	err = errors.Cause(err)
	//use type assertion to conver the err into XCode format
	if code, ok := err.(XCode); ok {
		// if type assertion success , call func convert the Code into grpc status
		grpcStatus, err := gRPCStatusFromXCode(code)
		if err == nil {
			return grpcStatus
		}
	}
	var grpcStatus *status.Status
	switch err {
	case context.Canceled:
		grpcStatus, _ = gRPCStatusFromXCode(Canceled)
	case context.DeadlineExceeded:
		grpcStatus, _ = gRPCStatusFromXCode(Deadline)
	default:
		if err == nil {
			return nil
		}
		// code := NewCode(123456, err.Error())
		// grpcStatus, _ = gRPCStatusFromXCode(code)
		grpcStatus, _ = status.FromError(err)
	}
	return grpcStatus
}

// gRPCStatusFromXCode make the Xcode convert into grpc status
func gRPCStatusFromXCode(code XCode) (*status.Status, error) {
	var sts *Status
	switch v := code.(type) {
	//if type is *Status
	case *Status:
		sts = v
	//if type is *Code, generate a new *Status with Code(return a XCode) and assert XCode to *Status
	case *Code:
		sts = NewStatus(*v).(*Status)
	default:
		fmt.Println("now in grpc from xcode default")
		sts = NewStatus(Code{code: code.Code(), msg: code.Message()}).(*Status)
		for _, detail := range code.Details() {
			if msg, ok := detail.(proto.Message); ok {
				_, _ = sts.WithDetails(msg)
			}
		}
	}
	//创建gRPC使用的status 将自定义的Status中的types.Status放入gRPC所使用的status中
	stas := status.New(codes.Unknown, strconv.Itoa(sts.Code()))
	return stas.WithDetails(sts.Proto())
}

// WithDetails put the proto message into status's details
func (s *Status) WithDetails(msgs ...proto.Message) (*Status, error) {
	for _, msg := range msgs {
		anyMSg, err := anypb.New(msg)
		if err != nil {
			return s, err
		}
		s.sts.Details = append(s.sts.Details, anyMSg)
	}
	return s, nil
}

// get the types.Status in struct Status
func (s *Status) Proto() *types.Status {
	return s.sts
}

// client receive the grpc server error , first use status.FromErr convert the
// err(send by grpc Server) into grpc status , secondly convert the grpc status into XCode
// thirdly use applet-api's ErrHandler recognize the XCode error and response into the user client

// GrpcStatusToXCode take out the details(protomessage) in gstatus,and convert it into XCode
func GrpcStatusToXCode(gstatus *status.Status) XCode {
	//if there is some XCode format message in details,convert the details's message into XCode
	details := gstatus.Details()
	for i := len(details) - 1; i >= 0; i-- {
		detail := details[i]
		if pb, ok := detail.(proto.Message); ok {
			return FromProto(pb)
		}
	}
	// if there is no some proto message in the gstatus's details
	return toXCode(*gstatus)
}

// FromProto convert  proto message (in gstatus deatils) to XCode ()
func FromProto(pbMsg proto.Message) XCode {
	//将传入的pbMsg断言为*types.Status，均拥有 ProtoReflect() 方法
	msg, ok := pbMsg.(*types.Status)
	if ok {
		if len(msg.Message) == 0 || msg.Message == strconv.FormatInt(int64(msg.Code), 10) {
			return &Code{code: int(msg.Code)}
		}
		return &Status{sts: msg}
	}
	code, _ := ServerErr.(*Code)
	return Errorf(*code, "invalid proto message get %v", pbMsg)
}

func Errorf(code Code, format string, args ...interface{}) XCode {
	code.msg = fmt.Sprintf(format, args...)
	return NewStatus(code)
}

// toXCode convert the grpc status into Code
func toXCode(grpcStatus status.Status) XCode {
	grpcCode := grpcStatus.Code()
	switch grpcCode {
	case codes.OK:
		return OK
	case codes.InvalidArgument:
		return RequestErr
	case codes.NotFound:
		return NotFound
	case codes.PermissionDenied:
		return AccessDenied
	case codes.Unauthenticated:
		return Unauthorized
	case codes.ResourceExhausted:
		return LimitExceed
	case codes.Unimplemented:
		return MethodNotAllowed
	case codes.DeadlineExceeded:
		return Deadline
	case codes.Unavailable:
		return ServiceUnavailable
	case codes.Unknown:
		return String(grpcStatus.Message())
	}
	return ServerErr
}
