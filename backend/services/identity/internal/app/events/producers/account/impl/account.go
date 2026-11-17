package impl
import(
 "context"
 "encoding/json"
 "time"
 "github.com/google/uuid"
 coreevents "grovehub/libs/core/events"
 infrastructure "grovehub/services/infrastructure/client/messaging"
 "grovehub/services/identity/internal/app/entity"
 producer "grovehub/services/identity/internal/app/events/producers/account"
 changed "grovehub/services/identity/internal/app/events/producers/account/operations/changed"
)
type Producer struct{events infrastructure.EventStore}
func New(events infrastructure.EventStore)producer.ChangedProducer{return &Producer{events:events}}
func snapshot(a *entity.Account)*changed.Snapshot{
 if a==nil{return nil}
 tenant:="";if a.TenantID!=nil{tenant=*a.TenantID}
 return &changed.Snapshot{AccountID:a.AccountID,PersonGID:a.PersonGID,TenantID:tenant,Username:a.Username,Status:string(a.Status),Version:a.Version}
}
func(p *Producer)Enqueue(ctx context.Context,operation string,old,next *entity.Account)error{
 eventID:=uuid.NewString();id:=""
 if next!=nil{id=next.AccountID}else if old!=nil{id=old.AccountID}
 message:=coreevents.NewEnvelope(ctx,coreevents.Metadata{EventID:eventID,EventType:"identity.account.changed.v1",Version:1,AggregateID:id,OccurredAt:time.Now().UTC()},changed.Changed{Operation:operation,Old:snapshot(old),New:snapshot(next)})
 raw,err:=json.Marshal(message);if err!=nil{return err}
 return p.events.Enqueue(ctx,infrastructure.Event{EventID:eventID,Topic:"identity_account_changed",Payload:raw,Metadata:map[string]string{"event-id":eventID,"event-type":"identity.account.changed.v1","source":"identity"}})
}
