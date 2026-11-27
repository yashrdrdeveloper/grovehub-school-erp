package resetpassword

import (
	"context"
)

// Operation is the use-case contract; concrete dependencies live in impl/.
type Operation interface {
	Execute(ctx context.Context, grant string, password string) error
}
