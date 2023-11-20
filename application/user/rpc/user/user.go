// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package user

import (
	"context"

	"github.com/GGjahoon/IZH/application/user/rpc/service"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	LoginRequest     = service.LoginRequest
	LoginResponse    = service.LoginResponse
	RegisterRequest  = service.RegisterRequest
	RegisterResponse = service.RegisterResponse

	User interface {
		Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
		Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	}

	defaultUser struct {
		cli zrpc.Client
	}
)

func NewUser(cli zrpc.Client) User {
	return &defaultUser{
		cli: cli,
	}
}

func (m *defaultUser) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	client := service.NewUserClient(m.cli.Conn())
	return client.Register(ctx, in, opts...)
}

func (m *defaultUser) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	client := service.NewUserClient(m.cli.Conn())
	return client.Login(ctx, in, opts...)
}
