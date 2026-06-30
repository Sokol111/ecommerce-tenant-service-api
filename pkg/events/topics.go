package events

import (
	eventsv1 "github.com/Sokol111/ecommerce-tenant-service-api/gen/events/tenant/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// TopicMap maps proto message full names to their Kafka topics.
var TopicMap = map[protoreflect.FullName]string{
	(&eventsv1.TenantUpdatedEvent{}).ProtoReflect().Descriptor().FullName(): TopicTenantTenantEvents,
	(&eventsv1.TenantDeletedEvent{}).ProtoReflect().Descriptor().FullName(): TopicTenantTenantEvents,
}
