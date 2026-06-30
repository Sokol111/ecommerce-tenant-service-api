package consumer

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-commons/pkg/tenant"
	eventsv1 "github.com/Sokol111/ecommerce-tenant-service-api/gen/events/tenant/v1"
	"go.uber.org/zap"
)

type tenantEventHandler struct {
	lifecycle tenant.Lifecycle
	log       *zap.Logger
}

func newTenantEventHandler(lifecycle tenant.Lifecycle, log *zap.Logger) *tenantEventHandler {
	return &tenantEventHandler{lifecycle: lifecycle, log: log}
}

func (h *tenantEventHandler) HandleTenantUpdated(ctx context.Context, evt *eventsv1.TenantUpdatedEvent) error {
	if !evt.GetEnabled() {
		h.log.Info("tenant disabled, skipping migration", zap.String("tenant", evt.GetSlug()))
		return nil
	}

	h.log.Info("tenant updated, running migration", zap.String("tenant", evt.GetSlug()))

	if err := h.lifecycle.Create(ctx, evt.GetSlug()); err != nil {
		return fmt.Errorf("tenant create failed for %q: %w", evt.GetSlug(), err)
	}

	return nil
}

func (h *tenantEventHandler) HandleTenantDeleted(ctx context.Context, evt *eventsv1.TenantDeletedEvent) error {
	h.log.Info("tenant deleted, scheduling deferred cleanup", zap.String("tenant", evt.GetSlug()))

	if err := h.lifecycle.Delete(ctx, evt.GetSlug()); err != nil {
		return fmt.Errorf("tenant delete failed for %q: %w", evt.GetSlug(), err)
	}

	return nil
}
