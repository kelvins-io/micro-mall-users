package server

import (
	"context"
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_users_proto/users"
)

type MerchantsServer struct{}

func NewMerchantsServer() users.MerchantsServiceServer {
	return new(MerchantsServer)
}

func (m *MerchantsServer) MerchantsMaterial(ctx context.Context, req *users.MerchantsMaterialRequest) (*users.MerchantsMaterialResponse, error) {
	return &users.MerchantsMaterialResponse{}, nil
}

func (m *MerchantsServer) MerchantsMaterialAudit(ctx context.Context, req *users.MerchantsMaterialAuditRequest) (*users.MerchantsMaterialAuditResponse, error) {
	return &users.MerchantsMaterialAuditResponse{}, nil
}

func (m *MerchantsServer) GetMerchantsMaterial(ctx context.Context, req *users.GetMerchantsMaterialRequest) (*users.GetMerchantsMaterialResponse, error) {
	return &users.GetMerchantsMaterialResponse{}, nil
}

func (m *MerchantsServer) MerchantsAssociateShop(ctx context.Context, req *users.MerchantsAssociateShopRequest) (*users.MerchantsAssociateShopResponse, error) {
	return &users.MerchantsAssociateShopResponse{}, nil
}
