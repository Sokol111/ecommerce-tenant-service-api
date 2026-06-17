package client

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/fx"

	grpcclient "github.com/Sokol111/ecommerce-commons/pkg/grpc/client"
	tenantv1 "github.com/Sokol111/ecommerce-tenant-service-api/gen/connect/tenant/v1"
)

// Module wires a native gRPC client for TenantService.
// Configuration is read from koanf under key "tenant.grpc".
func Module() fx.Option {
	return fx.Module("tenant-grpc-client",
		fx.Provide(func(k *koanf.Koanf) (grpcclient.Config, error) {
			return grpcclient.LoadConfig(k, "tenant.grpc")
		}, fx.Private),
		fx.Provide(grpcclient.NewConn, fx.Private),
		fx.Provide(tenantv1.NewTenantServiceClient),
	)
}
