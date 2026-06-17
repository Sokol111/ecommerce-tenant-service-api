// Package tenantevents provides a Kafka consumer module for tenant lifecycle events.
// It handles tenant migrations on create/update and deferred cleanup on delete.
package consumer

import (
	"github.com/Sokol111/ecommerce-commons/pkg/messaging/kafka/consumer"
	"github.com/Sokol111/ecommerce-tenant-service-api/gen/events"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const tenantEventsConsumer = "tenant-events"

// Module registers a Kafka consumer for tenant events that runs
// database migrations when tenants are created/updated and schedules deferred cleanup on delete.
//
// Requires tenant.NewModule() to be registered in the fx container for Registry and MigrationRunner.
// To add service-specific cleanup, register additional tenant.Cleaner implementations
// in the "tenant_cleaners" fx group:
//
//	fx.Provide(fx.Annotate(s3.NewImageTenantCleaner,
//	    fx.ResultTags(`group:"tenant_cleaners"`),
//	))
func Module() fx.Option {
	return fx.Options(
		events.Module(),
		fx.Provide(newTenantEventHandler),
		consumer.RegisterHandlerAndConsumer(tenantEventsConsumer, newTenantRouter),
	)
}

func newTenantRouter(h *tenantEventHandler, log *zap.Logger) consumer.Handler {
	r := consumer.NewRouter(log)
	consumer.Register(r, h.HandleTenantUpdated)
	consumer.Register(r, h.HandleTenantDeleted)
	return r
}
