package strategies

import "context"

// @Strategy
type Strategy interface { Execute(ctx context.Context) error }
