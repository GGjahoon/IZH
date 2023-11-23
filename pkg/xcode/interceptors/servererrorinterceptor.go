package interceptors

import (
	"context"

	"github.com/GGjahoon/IZH/pkg/xcode"
	"google.golang.org/grpc"
)

func ServerErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		resp, err = handler(ctx, req)
		return resp, xcode.FromErr(err).Err()
	}
}
