package messaging

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
