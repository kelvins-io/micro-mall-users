package client

import (
	"context"
	"gitee.com/cristiane/micro-mall-users/pkg/util"
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_users_proto/users"
	"testing"
)

const ServerName = "micro-mall-users"

func TestRegister(t *testing.T) {
	conn, err := util.GetGrpcClient(ServerName)
	if err != nil {
		t.Errorf("Conn err: %v", err)
	}
	defer conn.Close()

	client := users.NewUsersServiceClient(conn)
	ctx := context.Background()
	r := users.GetUserInfoRequest{
		Uid:                  10009,
	}
	accountInfo, err := client.GetUserInfo(ctx, &r)
	if err != nil {
		t.Error("GetInfoByAccountId err ", err)
	} else {
		t.Logf("accountInfo: %+v", accountInfo.Info)
	}
}
