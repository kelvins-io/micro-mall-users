package client

import (
	"context"
	"fmt"
	"gitee.com/cristiane/micro-mall-users/model/args"
	"gitee.com/cristiane/micro-mall-users/pkg/util"
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_users_proto/users"
	"gitee.com/kelvins-io/kelvins/util/test_tool"
	"github.com/bojand/ghz/runner"
	"os"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := util.GetGrpcClient(ctx, args.RpcServiceMicroMallUsers)
	if err != nil {
		t.Errorf("Conn err: %v", err)
	}
	//defer conn.Close()

	client := users.NewUsersServiceClient(conn)
	r := users.GetUserInfoRequest{
		Uid: 10009,
	}
	accountInfo, err := client.GetUserInfo(ctx, &r)
	if err != nil {
		t.Error("GetInfoByAccountId err ", err)
	} else {
		t.Logf("accountInfo: %+v", accountInfo.Info)
	}
}

type ghzReportConf struct {
	call     string
	dataJson string
}

func TestGHZ(t *testing.T) {
	testCases := []ghzReportConf{
		{
			call:     "GetUserInfo",
			dataJson: `{"uid":10048}`,
		},
		{
			call:     "GetUserInfoByPhone",
			dataJson: `{"country_code":"86","phone":"11606450640"}`,
		},
		{
			call:     "CheckUserByPhone",
			dataJson: `{"country_code":"86","phone":"11606450640"}`,
		},
		{
			call:     "CheckUserState",
			dataJson: `{"uid_list":[10049,10050,10051,10052]}`,
		},
	}
	var err error
	for _, v := range testCases {
		t.Run(v.call, func(t *testing.T) {
			err = testUsersServiceGhz(v)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func testUsersServiceGhz(conf ghzReportConf) error {
	gopath := os.Getenv("GOPATH")
	reportFile, err := os.Create(gopath + fmt.Sprintf("/src/gitee.com/cristiane/micro-mall-users/test_report/UsersService.%s.ghz.html", conf.call))
	if err != nil {
		return err
	}
	var opts []runner.Option
	opts = append(opts,
		runner.WithProtoFile(gopath+"/src/gitee.com/cristiane/micro-mall-users-proto/users/users.proto", []string{gopath + "/src/",}),
		runner.WithConcurrency(200),
		runner.WithConnections(5),
		runner.WithDataFromJSON(conf.dataJson),
		runner.WithInsecure(true),
		runner.WithTotalRequests(3000))

	return test_tool.ExecuteRPCGhzTest(&test_tool.GhzTestOption{
		Call:         "users.UsersService." + conf.call,
		Host:         "localhost:54786",
		Token:        "c9VW6ForlmzdeDkZE2i8",
		ReportFormat: test_tool.ReportHTML,
		Out:          reportFile,
		Options:      opts,
	})
}
