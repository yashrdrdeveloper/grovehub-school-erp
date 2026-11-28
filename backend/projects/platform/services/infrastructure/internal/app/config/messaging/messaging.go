package messaging

import (
    "errors"
    "fmt"
    "net/url"
    "os"
    "strings"
    "time"
)

const ProviderNATS = "nats"

type Messaging struct {
    Provider string
    NATSURL string
    QueueGroupPrefix string
    DurablePrefix string
    CloseTimeout time.Duration
    AckWaitTimeout time.Duration
    SubscribeTimeout time.Duration
}

// FromEnv loads runtime transport configuration without any fixed broker
// address or durable/group names embedded in Go application defaults.
func FromEnv() Messaging {
    provider := strings.TrimSpace(os.Getenv("GROVEHUB_MESSAGING_PROVIDER"))
    if provider == "" { provider = ProviderNATS }
    return Messaging{
        Provider: provider,
        NATSURL: strings.TrimSpace(os.Getenv("GROVEHUB_NATS_URL")),
        QueueGroupPrefix: strings.TrimSpace(os.Getenv("GROVEHUB_NATS_QUEUE_GROUP_PREFIX")),
        DurablePrefix: strings.TrimSpace(os.Getenv("GROVEHUB_NATS_DURABLE_PREFIX")),
        CloseTimeout: 30*time.Second,
        AckWaitTimeout: 30*time.Second,
        SubscribeTimeout: 30*time.Second,
    }
}

func (config Messaging) Validate() error {
    if strings.TrimSpace(config.Provider)=="" { return errors.New("GROVEHUB_MESSAGING_PROVIDER is required") }
    if config.Provider != ProviderNATS { return nil } // Registry validates supported providers.
    if config.NATSURL=="" { return errors.New("GROVEHUB_NATS_URL is required") }
    parsed,err:=url.Parse(config.NATSURL)
    if err!=nil || parsed.Host=="" || parsed.Scheme=="" {
        return fmt.Errorf("invalid GROVEHUB_NATS_URL: must be an absolute broker URL")
    }
    if config.QueueGroupPrefix=="" { return errors.New("GROVEHUB_NATS_QUEUE_GROUP_PREFIX is required") }
    if config.DurablePrefix=="" { return errors.New("GROVEHUB_NATS_DURABLE_PREFIX is required") }
    return nil
}

func DefaultMessaging() Messaging { return FromEnv() }
