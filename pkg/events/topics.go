package events

import (
	eventsv1 "github.com/Sokol111/ecommerce-tenant-service-api/gen/events/tenant/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Topic constants
const (
	TopicTenantTenantEvents = "tenant.tenant.events"
)

// topicMap maps proto message full names to their Kafka topics.
var topicMap = map[protoreflect.FullName]string{
	(&eventsv1.TenantUpdatedEvent{}).ProtoReflect().Descriptor().FullName(): TopicTenantTenantEvents,
	(&eventsv1.TenantDeletedEvent{}).ProtoReflect().Descriptor().FullName(): TopicTenantTenantEvents,
}

func init() {
	fd := eventsv1.File_tenant_v1_events_proto
	msgs := fd.Messages()
	for i := range msgs.Len() {
		fullName := msgs.Get(i).FullName()
		if _, ok := topicMap[fullName]; !ok {
			panic("events: no topic registered for " + string(fullName))
		}
	}
}

// TopicFor returns the Kafka topic for the given proto message.
// Panics if the message type is not registered in topicMap.
func TopicFor(msg proto.Message) string {
	fullName := msg.ProtoReflect().Descriptor().FullName()
	topic, ok := topicMap[fullName]
	if !ok {
		panic("events: no topic registered for " + string(fullName))
	}
	return topic
}
