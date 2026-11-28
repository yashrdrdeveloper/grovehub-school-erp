// Package monitoring tracks process-local Audit readiness and bounded metrics.
package monitoring

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	coretelemetry "grovehub/libs/core/telemetry"
)

const ProbeInterval = 10 * time.Second
const ProbeTimeout = 2 * time.Second
const ProbeMaxAge = 30 * time.Second

type topicState struct {
	subscribed bool
	failed     bool
}

type Monitor struct {
	mu                sync.RWMutex
	topics            map[string]topicState
	probedAt          time.Time
	dependenciesReady bool
	closed            bool
	deliveries        metric.Int64Counter
	duration          metric.Float64Histogram
	reconnects        metric.Int64Counter
	registration      metric.Registration
}

func New(topics map[string]string) (*Monitor, error) {
	monitor := &Monitor{topics: make(map[string]topicState, len(topics))}
	for topic := range topics {
		monitor.topics[topic] = topicState{}
	}
	meter := coretelemetry.MeterProvider(coretelemetry.ServiceAudit.Key).Meter("grovehub/audit/consumer")
	var err error
	monitor.deliveries, err = meter.Int64Counter("grovehub.audit.consumer.deliveries", metric.WithUnit("{event}"))
	if err != nil {
		return nil, err
	}
	monitor.duration, err = meter.Float64Histogram("grovehub.audit.consumer.duration", metric.WithUnit("s"))
	if err != nil {
		return nil, err
	}
	monitor.reconnects, err = meter.Int64Counter("grovehub.audit.consumer.interruptions", metric.WithUnit("{interruption}"))
	if err != nil {
		return nil, err
	}
	gauge, err := meter.Int64ObservableGauge("grovehub.audit.consumer.subscribed", metric.WithUnit("1"))
	if err != nil {
		return nil, err
	}
	ready, err := meter.Int64ObservableGauge("grovehub.audit.ready", metric.WithUnit("1"))
	if err != nil {
		return nil, err
	}
	monitor.registration, err = meter.RegisterCallback(func(_ context.Context, observer metric.Observer) error {
		monitor.mu.RLock()
		defer monitor.mu.RUnlock()
		for topic, state := range monitor.topics {
			observer.ObserveInt64(gauge, boolValue(state.subscribed && !monitor.closed), metric.WithAttributes(attribute.String("messaging.destination.name", topic)))
		}
		observer.ObserveInt64(ready, boolValue(monitor.readyLocked(time.Now())))
		return nil
	}, gauge, ready)
	if err != nil {
		return nil, err
	}
	return monitor, nil
}

func boolValue(value bool) int64 {
	if value {
		return 1
	}
	return 0
}

func (monitor *Monitor) Subscription(topic string, ready bool) {
	monitor.mu.Lock()
	defer monitor.mu.Unlock()
	if state, ok := monitor.topics[topic]; ok && !monitor.closed {
		state.subscribed = ready
		monitor.topics[topic] = state
	}
}

func (monitor *Monitor) Delivery(ctx context.Context, topic string, outcome string, elapsed time.Duration, err error) {
	monitor.mu.Lock()
	state, known := monitor.topics[topic]
	if known && !monitor.closed {
		state.failed = err != nil
		monitor.topics[topic] = state
	}
	monitor.mu.Unlock()
	if !known {
		return
	}
	if err != nil {
		outcome = "retry"
	}
	switch outcome {
	case "accepted", "invalid_event", "identity_conflict", "retry":
	default:
		outcome = "retry"
	}
	attrs := metric.WithAttributes(attribute.String("messaging.destination.name", topic), attribute.String("grovehub.audit.outcome", outcome))
	monitor.deliveries.Add(ctx, 1, attrs)
	monitor.duration.Record(ctx, elapsed.Seconds(), attrs)
}

func (monitor *Monitor) Interrupted(ctx context.Context, topic string) {
	monitor.Subscription(topic, false)
	monitor.mu.RLock()
	_, known := monitor.topics[topic]
	monitor.mu.RUnlock()
	if known {
		monitor.reconnects.Add(ctx, 1, metric.WithAttributes(attribute.String("messaging.destination.name", topic)))
	}
}

func (monitor *Monitor) Probe(at time.Time, ready bool) {
	monitor.mu.Lock()
	defer monitor.mu.Unlock()
	monitor.probedAt = at
	monitor.dependenciesReady = ready
}

func (monitor *Monitor) Ready() bool {
	monitor.mu.RLock()
	defer monitor.mu.RUnlock()
	return monitor.readyLocked(time.Now())
}

func (monitor *Monitor) readyLocked(now time.Time) bool {
	if monitor.closed || !monitor.dependenciesReady || monitor.probedAt.IsZero() || now.Sub(monitor.probedAt) > ProbeMaxAge || len(monitor.topics) == 0 {
		return false
	}
	for _, state := range monitor.topics {
		if !state.subscribed || state.failed {
			return false
		}
	}
	return true
}

func (monitor *Monitor) Close() {
	monitor.mu.Lock()
	if monitor.closed {
		monitor.mu.Unlock()
		return
	}
	monitor.closed = true
	monitor.mu.Unlock()
	if monitor.registration != nil {
		_ = monitor.registration.Unregister()
	}
}
