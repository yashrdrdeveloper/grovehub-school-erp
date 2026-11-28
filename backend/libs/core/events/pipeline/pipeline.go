// Package pipeline exposes transport-agnostic, typed event eligibility and delivery rules.
package pipeline

import (
	"context"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/internal/src/operations/pipeline"
)

var ErrIncompleteRules = contracts.ErrIncompleteRules

// Rules composes eligibility, construction and delivery using a consumer-owned message.
type Rules[T any, M any] = contracts.Rules[T, M]

func Run[T any, M any](ctx context.Context, value T, rules Rules[T, M]) error {
	return pipeline.Run(ctx, value, rules)
}
