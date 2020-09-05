package server

import (
	"context"
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_users_proto/users"
)

type UsersServer struct {
}

func NewUsersServer() users.UsersServiceServer {
	return new(UsersServer)
}

func (u *UsersServer) GetUserInfo(ctx context.Context, req *users.GetUserInfoRequest) (*users.GetUserInfoResponse, error) {
	return &users.GetUserInfoResponse{}, nil
}
