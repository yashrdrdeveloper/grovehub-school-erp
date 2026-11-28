package events_test

import (
    "context"
    "encoding/json"
    "time"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    "go.opentelemetry.io/otel/trace"

    coreevents "grovehub/libs/core/events"
)

var _ = Describe("Generic event envelope", func() {
    type payload struct {
        GID string `json:"gid"`
    }

    It("serializes a generic domain payload with all defined metadata", func() {
        occurredAt := time.Date(2026, 11, 6, 10, 0, 0, 0, time.UTC)
        ctx := coreevents.WithCausationID(coreevents.WithCorrelationID(context.Background(), "req_789"), "cmd_456")
        event := coreevents.NewEnvelope(ctx, coreevents.Metadata{
            EventID: "evt_1001", EventType: "account.created", Version: 1,
            AggregateID: "account_123", OccurredAt: occurredAt,
        }, payload{GID: "person_123"})

        raw, err := json.Marshal(event)
        Expect(err).NotTo(HaveOccurred())
        var fields map[string]any
        Expect(json.Unmarshal(raw, &fields)).To(Succeed())
        Expect(fields).To(HaveLen(8))
        Expect(fields["eventId"]).To(Equal("evt_1001"))
        Expect(fields["eventType"]).To(Equal("account.created"))
        Expect(fields["version"]).To(BeNumerically("==", 1))
        Expect(fields["aggregateId"]).To(Equal("account_123"))
        Expect(fields["occurredAt"]).To(Equal("2026-11-06T10:00:00Z"))
        Expect(fields["correlationId"]).To(Equal("req_789"))
        Expect(fields["causationId"]).To(Equal("cmd_456"))
        Expect(fields["data"]).To(Equal(map[string]any{"gid": "person_123"}))
        Expect(fields).NotTo(HaveKey("source"))

        var roundTrip coreevents.Envelope[payload]
        Expect(json.Unmarshal(raw, &roundTrip)).To(Succeed())
        Expect(roundTrip.Data.GID).To(Equal("person_123"))
    })

    It("uses the OpenTelemetry trace ID when no correlation ID was provided", func() {
        traceID, err := trace.TraceIDFromHex("0123456789abcdef0123456789abcdef")
        Expect(err).NotTo(HaveOccurred())
        spanID, err := trace.SpanIDFromHex("0123456789abcdef")
        Expect(err).NotTo(HaveOccurred())
        ctx := trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
            TraceID: traceID, SpanID: spanID, TraceFlags: trace.FlagsSampled,
        }))
        event := coreevents.NewEnvelope(ctx, coreevents.Metadata{
            EventID: "e", EventType: "person.changed", Version: 1, AggregateID: "p", OccurredAt: time.Now().UTC(),
        }, payload{GID: "p"})
        Expect(event.CorrelationID).To(Equal(traceID.String()))
        Expect(event.CausationID).To(BeEmpty())
    })

    It("omits genuinely unknown correlation and causation IDs", func() {
        event := coreevents.NewEnvelope(context.Background(), coreevents.Metadata{
            EventID: "e", EventType: "person.changed", Version: 1, AggregateID: "p", OccurredAt: time.Now().UTC(),
        }, payload{GID: "p"})
        raw, err := json.Marshal(event)
        Expect(err).NotTo(HaveOccurred())
        var fields map[string]any
        Expect(json.Unmarshal(raw, &fields)).To(Succeed())
        Expect(fields).NotTo(HaveKey("correlationId"))
        Expect(fields).NotTo(HaveKey("causationId"))
    })

    It("prefers explicit metadata IDs over context IDs", func() {
        ctx := coreevents.WithCausationID(coreevents.WithCorrelationID(context.Background(), "req_ctx"), "cmd_ctx")
        event := coreevents.NewEnvelope(ctx, coreevents.Metadata{
            EventID: "e", EventType: "x", Version: 1, AggregateID: "a",
            CorrelationID: "req_metadata", CausationID: "cmd_metadata",
        }, payload{})
        Expect(event.CorrelationID).To(Equal("req_metadata"))
        Expect(event.CausationID).To(Equal("cmd_metadata"))
    })
})
