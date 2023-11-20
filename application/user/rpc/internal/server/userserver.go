// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package server

import (
	"context"

	"github.com/GGjahoon/IZH/application/user/rpc/internal/logic"
	"github.com/GGjahoon/IZH/application/user/rpc/internal/svc"
	"github.com/GGjahoon/IZH/application/user/rpc/service"
)

type UserServer struct {
	svcCtx *svc.ServiceContext
	service.UnimplementedUserServer
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{
		svcCtx: svcCtx,
	}
}

func (s *UserServer) Register(ctx context.Context, in *service.RegisterRequest) (*service.RegisterResponse, error) {
	l := logic.NewRegisterLogic(ctx, s.svcCtx)
	return l.Register(in)
}

func (s *UserServer) Login(ctx context.Context, in *service.LoginRequest) (*service.LoginResponse, error) {
	l := logic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(in)
}
