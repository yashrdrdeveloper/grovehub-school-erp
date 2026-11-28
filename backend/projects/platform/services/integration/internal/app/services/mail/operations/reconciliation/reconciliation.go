package reconciliation

import (
	"context"
	"time"

	"grovehub/services/integration/client/mail"
	"grovehub/services/identity/client/invitation"
	"grovehub/services/integration/internal/app/repository/maildelivery"
)

const BatchSize = 256

// Snapshot executes a read-only operation against one consistent database snapshot.
type Snapshot func(context.Context, func(context.Context,time.Time) error) error

type Operation interface { mail.Client }

type JobReader = invitation.DeliveryReader
type ReceiptReader = maildelivery.ReceiptReader
