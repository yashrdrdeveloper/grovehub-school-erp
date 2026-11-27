package messaging

import (
    "strings"
    "testing"
)

func TestFromEnvUsesConfiguredBroker(t *testing.T) {
    t.Setenv("GROVEHUB_MESSAGING_PROVIDER", "nats")
    t.Setenv("GROVEHUB_NATS_URL", "nats://broker.internal:5222")
    t.Setenv("GROVEHUB_NATS_QUEUE_GROUP_PREFIX", "school-workers")
    t.Setenv("GROVEHUB_NATS_DURABLE_PREFIX", "school-durable")
    config:=FromEnv()
    if config.NATSURL!="nats://broker.internal:5222" || config.QueueGroupPrefix!="school-workers" || config.DurablePrefix!="school-durable" {t.Fatalf("invalid env configuration: %+v",config)}
    if err:=config.Validate();err!=nil {t.Fatal(err)}
}

func TestFromEnvRequiresBrokerSettings(t *testing.T) {
    t.Setenv("GROVEHUB_MESSAGING_PROVIDER", "nats")
    t.Setenv("GROVEHUB_NATS_URL","")
    t.Setenv("GROVEHUB_NATS_QUEUE_GROUP_PREFIX","school")
    t.Setenv("GROVEHUB_NATS_DURABLE_PREFIX","school")
    if err:=FromEnv().Validate();err==nil||!strings.Contains(err.Error(),"GROVEHUB_NATS_URL"){t.Fatalf("missing broker URL should fail: %v",err)}
}
