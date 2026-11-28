package pipeline

import (
	"context"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/internal/src/contracts"
)

func Run[T any, M any](ctx context.Context, value T, rules contracts.Rules[T, M]) error {
	if rules.ShouldEmit == nil || rules.Build == nil || rules.Enqueue == nil {
		return contracts.ErrIncompleteRules
	}
	if !rules.ShouldEmit(value) {
		return nil
	}
	message, err := rules.Build(ctx, value)
	if err != nil {
		return err
	}
	return rules.Enqueue(ctx, message)
}
