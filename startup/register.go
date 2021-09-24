package startup

import (
	"context"
	"gitee.com/kelvins-io/kelvins"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net/http"
	"time"

	"gitee.com/cristiane/micro-mall-users/http_server"
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_users_proto/users"
	"gitee.com/cristiane/micro-mall-users/server"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// RegisterGRPCServer 此处注册pb的Server
func RegisterGRPCServer(grpcServer *grpc.Server) error {
	users.RegisterUsersServiceServer(grpcServer, server.NewUsersServer())
	users.RegisterMerchantsServiceServer(grpcServer, server.NewMerchantsServer())
	return nil
}

// RegisterGateway 此处注册pb的Gateway
func RegisterGateway(ctx context.Context, gateway *runtime.ServeMux, endPoint string, dopts []grpc.DialOption) error {
	if err := users.RegisterUsersServiceHandlerFromEndpoint(ctx, gateway, endPoint, dopts); err != nil {
		return err
	}
	if err := users.RegisterMerchantsServiceHandlerFromEndpoint(ctx, gateway, endPoint, dopts); err != nil {
		return err
	}
	return nil
}

// RegisterHttpRoute 此处注册http接口
func RegisterHttpRoute(serverMux *http.ServeMux) error {
	serverMux.HandleFunc("/swagger/", http_server.SwaggerHandler)
	return nil
}

const (
	healthCheckSleep = 30 * time.Second
)

func RegisterGRPCHealthStatusHandle(health *kelvins.GRPCHealthServer) {
	// first set
	health.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	health.SetServingStatus("users.UsersService", healthpb.HealthCheckResponse_SERVING)
	health.SetServingStatus("users.MerchantsService", healthpb.HealthCheckResponse_SERVING)
	ticker := time.NewTicker(healthCheckSleep)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			health.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
			health.SetServingStatus("users.UsersService", healthpb.HealthCheckResponse_SERVING)
			health.SetServingStatus("users.MerchantsService", healthpb.HealthCheckResponse_SERVING)
		case <-kelvins.AppCloseCh:
			health.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
			health.SetServingStatus("users.UsersService", healthpb.HealthCheckResponse_NOT_SERVING)
			health.SetServingStatus("users.MerchantsService", healthpb.HealthCheckResponse_NOT_SERVING)
			return
		}
	}
}
