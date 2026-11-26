package activity

import (
	"context"
	"grovehub/services/audit/internal/app/events/consumers/activity/helpers"
	"grovehub/services/audit/internal/app/repository/activity"
	"grovehub/services/infrastructure/client/messaging"
)

type Consumer struct { repository activity.Repository }

func New(repository activity.Repository) *Consumer {
	return &Consumer{repository:repository}
}

func (consumer *Consumer) Handle(ctx context.Context, eventType string, message messaging.Message) error {
	entry, err := helpers.Project(ctx,eventType,message.Payload)
	if err != nil { return err }
	return consumer.repository.Insert(ctx,entry)
}
