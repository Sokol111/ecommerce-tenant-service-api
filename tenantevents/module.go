// Package tenantevents provides a Kafka consumer module for tenant lifecycle events.
// It handles tenant migrations on create/update and cleanup on delete.
package tenantevents

import (
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
	tenant_events "github.com/Sokol111/ecommerce-tenant-service-api/gen/events"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const tenantEventsConsumer = "tenant-events"

// Module registers a Kafka consumer for tenant events that runs
// database migrations when tenants are created/updated and invokes cleaners on delete.
//
// The MongoDB database cleanup is always included.
// To add service-specific cleanup, register additional tenant.Cleaner implementations
// in the "tenant_cleaners" fx group:
//
//	fx.Provide(fx.Annotate(s3.NewImageTenantCleaner,
//	    fx.ResultTags(`group:"tenant_cleaners"`),
//	))
func Module() fx.Option {
	return fx.Options(
		tenant_events.Module(),
		fx.Provide(
			fx.Annotate(
				commonsmongo.NewTenantCleanupCleaner,
				fx.ResultTags(`group:"tenant_cleaners"`),
			),
			newTenantEventHandler,
		),
		consumer.RegisterHandlerAndConsumer(tenantEventsConsumer, newTenantRouter),
	)
}

func newTenantRouter(h *tenantEventHandler, log *zap.Logger) consumer.Handler {
	r := consumer.NewRouter(log)
	consumer.Register(r, h.HandleTenantUpdated)
	consumer.Register(r, h.HandleTenantDeleted)
	return r
}
