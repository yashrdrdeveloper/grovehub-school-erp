package delivery

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	client "grovehub/services/integration/client/connection"
	"grovehub/services/integration/internal/app/repository/maildelivery"
)

type Dispatcher struct { repository maildelivery.Repository }

func New(repository maildelivery.Repository) *Dispatcher {
	return &Dispatcher{repository: repository}
}

func digest(input []byte) string {
	sum := sha256.Sum256(input)
	return hex.EncodeToString(sum[:])
}

func (dispatcher *Dispatcher) Replay(ctx context.Context, tenantID string, deliveryID string, input []byte) (bool, error) {
	record, found, err := dispatcher.repository.Find(ctx, tenantID, deliveryID)
	if err != nil || !found { return found, err }
	if record.RequestHash != digest(input) { return true, client.ErrMailDeliveryConflict }
	if record.State != maildelivery.Accepted { return true, client.ErrMailDeliveryUncertain }
	return true, nil
}

func (dispatcher *Dispatcher) Dispatch(ctx context.Context, tenantID string, deliveryID string, input []byte, send func() (any, error)) (any, error) {
	if tenantID == "" || deliveryID == "" || send == nil { return nil, errors.New("invalid mail dispatch") }
	hash := digest(input)
	reserved, err := dispatcher.repository.Reserve(ctx, tenantID, deliveryID, hash)
	if err != nil { return nil, err }
	if !reserved {
		found, err := dispatcher.Replay(ctx, tenantID, deliveryID, input)
		if err != nil { return nil, err }
		if !found { return nil, client.ErrMailDeliveryUncertain }
		return nil, nil
	}
	result, err := send()
	if err != nil {
		// A timeout or transport error does not prove that the provider did not accept the email.
		return nil, client.ErrMailDeliveryUncertain
	}
	// Caller cancellation after provider acceptance must not undo the durable receipt.
	receiptCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 3*time.Second)
	defer cancel()
	if err := dispatcher.repository.Accept(receiptCtx, tenantID, deliveryID, hash); err != nil {
		return nil, client.ErrMailDeliveryUncertain
	}
	return result, nil
}
