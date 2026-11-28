package helpers

import "context"

type consumerGroupKey struct{}

// WithConsumerGroup selects an independent logical subscriber. Only trusted
// service wiring sets this value; never derive it from HTTP input.
func WithConsumerGroup(ctx context.Context, group string) context.Context {
	return context.WithValue(ctx,consumerGroupKey{},group)
}

func ConsumerGroup(ctx context.Context) string {
	group,_:=ctx.Value(consumerGroupKey{}).(string)
	return group
}

// SubscriptionObserver reports successful subscription and termination, not
// delivery of the first event. The callback must be safe for concurrent use.
type subscriptionObserverKey struct{}

func WithSubscriptionObserver(ctx context.Context, observer func(bool)) context.Context {
	return context.WithValue(ctx, subscriptionObserverKey{}, observer)
}

func NotifySubscription(ctx context.Context, ready bool) {
	if observer,ok:=ctx.Value(subscriptionObserverKey{}).(func(bool));ok && observer!=nil { observer(ready) }
}
