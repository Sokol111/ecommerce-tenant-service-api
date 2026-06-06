package tenantevents

import (
	"context"
	"fmt"

	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
	"github.com/Sokol111/ecommerce-commons/pkg/tenant"
	tenant_events "github.com/Sokol111/ecommerce-tenant-service-api/gen/events"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type tenantEventHandler struct {
	migrator commonsmongo.TenantMigrationRunner
	cleaners []tenant.Cleaner
	log      *zap.Logger
}

type tenantEventHandlerParams struct {
	fx.In

	Migrator commonsmongo.TenantMigrationRunner
	Cleaners []tenant.Cleaner `group:"tenant_cleaners"`
	Log      *zap.Logger
}

func newTenantEventHandler(p tenantEventHandlerParams) *tenantEventHandler {
	return &tenantEventHandler{migrator: p.Migrator, cleaners: p.Cleaners, log: p.Log}
}

func (h *tenantEventHandler) HandleTenantUpdated(ctx context.Context, evt *tenant_events.TenantUpdatedEvent) error {
	if !evt.Payload.Enabled {
		h.log.Info("tenant disabled, skipping migration", zap.String("tenant", evt.Payload.Slug))
		return nil
	}

	h.log.Info("tenant updated, running migration", zap.String("tenant", evt.Payload.Slug))

	if err := h.migrator.MigrateTenant(ctx, evt.Payload.Slug); err != nil {
		return fmt.Errorf("failed to migrate tenant %q: %w", evt.Payload.Slug, err)
	}

	return nil
}

func (h *tenantEventHandler) HandleTenantDeleted(ctx context.Context, evt *tenant_events.TenantDeletedEvent) error {
	h.log.Info("tenant deleted, cleaning up", zap.String("tenant", evt.Payload.Slug))

	for _, cleaner := range h.cleaners {
		if err := cleaner.CleanupTenant(ctx, evt.Payload.Slug); err != nil {
			return fmt.Errorf("failed to cleanup tenant %q: %w", evt.Payload.Slug, err)
		}
	}

	return nil
}
