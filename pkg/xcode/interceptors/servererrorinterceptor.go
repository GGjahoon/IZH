package interceptors

import (
	"context"
	"fmt"

	"github.com/GGjahoon/IZH/pkg/xcode"
	"google.golang.org/grpc"
)

func ServerErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		resp, err = handler(ctx, req)
		fmt.Println("handler over")
		xxx := xcode.FromErr(err).Err()
		fmt.Println(xxx)
		return resp, xxx
	}
}
