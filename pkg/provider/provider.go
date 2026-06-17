package provider

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-commons/pkg/tenant"
	tenantv1 "github.com/Sokol111/ecommerce-tenant-service-api/gen/connect/tenant/v1"
	"go.uber.org/fx"
)

type tenantSlugsProvider struct {
	client tenantv1.TenantServiceClient
}

func newTenantSlugsProvider(client tenantv1.TenantServiceClient) tenant.SlugsProvider {
	return &tenantSlugsProvider{client: client}
}

func (p *tenantSlugsProvider) GetSlugs(ctx context.Context) ([]string, error) {
	res, err := p.client.GetEnabledTenantSlugs(ctx, &tenantv1.GetEnabledTenantSlugsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tenant slugs: %w", err)
	}

	return res.Slugs, nil
}

// Module provides TenantSlugsProvider backed by the tenant-service API.
func Module() fx.Option {
	return fx.Module("tenant-slugs",
		fx.Provide(newTenantSlugsProvider),
	)
}
