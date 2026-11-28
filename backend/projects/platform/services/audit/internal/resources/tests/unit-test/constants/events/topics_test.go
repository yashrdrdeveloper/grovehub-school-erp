package events_test

import (
	"testing"

	"grovehub/services/audit/internal/app/constants/events"
)

func TestTopics(t *testing.T) {
	topics := events.Topics()
	if len(topics) != 9 {
		t.Fatalf("unexpected allowlist size: %d", len(topics))
	}
	for topic, kind := range topics {
		got, ok := events.Topic(kind)
		if !ok || got != topic {
			t.Fatalf("reverse lookup: %s -> %s", kind, got)
		}
	}
	delete(topics, events.ActivityOccurredTopic)
	if _, ok := events.Topic(events.ActivityOccurred); !ok {
		t.Fatal("caller modified the allowlist")
	}
	if _, ok := events.Topic("unregistered"); ok {
		t.Fatal("unknown event type accepted")
	}
}
