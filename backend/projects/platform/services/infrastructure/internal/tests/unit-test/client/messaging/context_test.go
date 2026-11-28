package messaging_test

import (
	"context"
	"reflect"
	"testing"

	"grovehub/services/infrastructure/client/messaging"
)

func TestConsumerGroup(t *testing.T) {
	ctx := context.Background()
	if group := messaging.ConsumerGroup(ctx); group != "" {
		t.Fatalf("unexpected default group: %q", group)
	}
	first := messaging.WithConsumerGroup(ctx, "audit")
	second := messaging.WithConsumerGroup(first, "other")
	if messaging.ConsumerGroup(first) != "audit" || messaging.ConsumerGroup(second) != "other" || messaging.ConsumerGroup(ctx) != "" {
		t.Fatal("consumer group must be isolated to its derived context")
	}
}

func TestSubscriptionObserver(t *testing.T) {
	ctx := context.Background()
	messaging.NotifySubscription(ctx, true)
	messaging.NotifySubscription(messaging.WithSubscriptionObserver(ctx, nil), false)
	var observed []bool
	subscription := messaging.WithSubscriptionObserver(messaging.WithConsumerGroup(ctx, "audit"), func(ready bool) {
		observed = append(observed, ready)
	})
	messaging.NotifySubscription(subscription, true)
	messaging.NotifySubscription(subscription, false)
	if !reflect.DeepEqual(observed, []bool{true, false}) || messaging.ConsumerGroup(subscription) != "audit" {
		t.Fatalf("unexpected subscription context: %v", observed)
	}
}

func TestPublicConfig(t *testing.T) {
	t.Setenv("GROVEHUB_MESSAGING_PROVIDER", "")
	t.Setenv("GROVEHUB_NATS_URL", "nats://broker:4222")
	t.Setenv("GROVEHUB_NATS_QUEUE_GROUP_PREFIX", "school")
	t.Setenv("GROVEHUB_NATS_DURABLE_PREFIX", "school")
	config := messaging.DefaultMessaging()
	if config.Provider != messaging.ProviderNATS || config != messaging.FromEnv() {
		t.Fatal("public config must use Infrastructure configuration")
	}
	if err := config.Validate(); err != nil {
		t.Fatal(err)
	}
}
