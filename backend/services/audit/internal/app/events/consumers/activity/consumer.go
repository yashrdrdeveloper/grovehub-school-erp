package activity

import (
	"context"
	"errors"
	"grovehub/services/audit/internal/app/repository/rejection"
	"grovehub/services/audit/internal/app/events/consumers/activity/helpers"
	"grovehub/services/audit/internal/app/repository/activity"
	"grovehub/services/infrastructure/client/messaging"
)

type Consumer struct { repository activity.Repository; rejected rejection.Repository }

func New(repository activity.Repository, rejected rejection.Repository) *Consumer {
	return &Consumer{repository:repository,rejected:rejected}
}

func (consumer *Consumer) Handle(ctx context.Context, eventType string, message messaging.Message) error {
	entry, err := helpers.Project(ctx,eventType,message.Payload)
	if errors.Is(err,helpers.ErrInvalidEvent) {
		return consumer.reject(ctx,eventType,message.Payload,"INVALID_EVENT",err)
	}
	if err!=nil { return err }
	err=consumer.repository.Insert(ctx,entry)
	if errors.Is(err,activity.ErrConflict) {
		return consumer.reject(ctx,eventType,message.Payload,"IDENTITY_CONFLICT",err)
	}
	return err
}

func (consumer *Consumer) reject(ctx context.Context,eventType string,payload []byte,reason string,original error) error {
	if consumer.rejected==nil { return original }
	return consumer.rejected.Record(ctx,eventType,payload,reason)
}
