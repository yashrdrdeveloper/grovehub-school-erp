package impl

import (
    "context"
    "encoding/json"

    mailbox "grovehub/services/integration/internal/app/events/producers/connection/operations/mailbox"
    helpers "grovehub/services/integration/internal/app/events/producers/connection/operations/mailbox/helpers"
    strategies "grovehub/services/integration/internal/app/events/producers/connection/operations/mailbox/strategies"
)

type ChangedStrategy struct {}
var _ strategies.ChangedStrategy = (*ChangedStrategy)(nil)

func New() strategies.ChangedStrategy { return &ChangedStrategy{} }

func (strategy *ChangedStrategy) Build(ctx context.Context, payload mailbox.Payload) (string, []byte, map[string]string, error) {
    event:=helpers.BuildChanged(ctx,payload)
    encoded,err:=json.Marshal(event)
    if err!=nil {return "",nil,nil,err}
    return mailbox.ChangedTopic,encoded,helpers.ChangedMetadata(event),nil
}
