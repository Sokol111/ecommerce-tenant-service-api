package tenantevents

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-commons/pkg/tenant"
	tenant_events "github.com/Sokol111/ecommerce-tenant-service-api/gen/events"
	"go.uber.org/zap"
)

type tenantEventHandler struct {
	lifecycle tenant.Lifecycle
	log       *zap.Logger
}

func newTenantEventHandler(lifecycle tenant.Lifecycle, log *zap.Logger) *tenantEventHandler {
	return &tenantEventHandler{lifecycle: lifecycle, log: log}
}

func (h *tenantEventHandler) HandleTenantUpdated(ctx context.Context, evt *tenant_events.TenantUpdatedEvent) error {
	if !evt.Payload.Enabled {
		h.log.Info("tenant disabled, skipping migration", zap.String("tenant", evt.Payload.Slug))
		return nil
	}

	h.log.Info("tenant updated, running migration", zap.String("tenant", evt.Payload.Slug))

	if err := h.lifecycle.Create(ctx, evt.Payload.Slug); err != nil {
		return fmt.Errorf("tenant create failed for %q: %w", evt.Payload.Slug, err)
	}

	return nil
}

func (h *tenantEventHandler) HandleTenantDeleted(ctx context.Context, evt *tenant_events.TenantDeletedEvent) error {
	h.log.Info("tenant deleted, scheduling deferred cleanup", zap.String("tenant", evt.Payload.Slug))

	if err := h.lifecycle.Delete(ctx, evt.Payload.Slug); err != nil {
		return fmt.Errorf("tenant delete failed for %q: %w", evt.Payload.Slug, err)
	}

	return nil
}
