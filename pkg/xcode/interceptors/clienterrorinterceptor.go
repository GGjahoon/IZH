package interceptors

import (
	"context"
	"fmt"

	"github.com/GGjahoon/IZH/pkg/xcode"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func ClientErrorInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			grpcStatus, _ := status.FromError(err)
			xcode := xcode.GrpcStatusToXCode(grpcStatus)
			err = errors.WithMessage(xcode, grpcStatus.Message())
			fmt.Println(grpcStatus)
			fmt.Println(grpcStatus.Details()...)
		}
		return err
	}
}
