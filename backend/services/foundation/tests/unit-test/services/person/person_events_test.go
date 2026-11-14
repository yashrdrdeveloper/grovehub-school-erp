package person_test

import (
    "context"
    "errors"
    "testing"
    "time"

    annotationdto "grovehub/libs/core/annotation/dto"
    foundationevents "grovehub/services/foundation/constants/events"
    persondto "grovehub/services/foundation/dto/person"
    "grovehub/services/foundation/entity"
    personservice "grovehub/services/foundation/services/person/impl"
)

type recordedPersonChange struct {
    operation foundationevents.Operation
    old *persondto.Person
    next *persondto.Person
}

type fakePersonChangeProducer struct {
    changes []recordedPersonChange
    err error
}

func (producer *fakePersonChangeProducer) Enqueue(ctx context.Context, operation foundationevents.Operation, old, next *persondto.Person) error {
    if ctx.Value(transactionMarker{}) != true { return errors.New("Person event enqueue must join the business transaction") }
    if producer.err != nil { return producer.err }
    producer.changes = append(producer.changes, recordedPersonChange{operation: operation, old: old, next: next})
    return nil
}

func TestPersonCreateEmitsOnlyOnceInsideUnitOfWork(t *testing.T) {
    people, contacts, work, _ := setup(t)
    producer := &fakePersonChangeProducer{}
    service := personservice.NewPersonService(people, contacts, work, fakeMapper{}, fakeValidator{}, producer)
    request := persondto.CreatePersonRequest{RequestBody: persondto.CreatePerson{
        PersonNames: persondto.PersonNames{FirstName: "Anu", LastName: "Roy"},
        Contacts: persondto.Contacts{Emails: []persondto.ContactValue{{Value: "anu@example.com", IsPrimary: true}}},
    }}
    createdPerson, created, err := service.Create(context.Background(), request)
    if err != nil || !created { t.Fatalf("create failed: %v, %v",created,err) }
    if len(producer.changes)!=1 || producer.changes[0].operation!=foundationevents.Created ||
        producer.changes[0].old!=nil || producer.changes[0].next.GID!=createdPerson.GID {
        t.Fatalf("unexpected Create event: %#v",producer.changes)
    }
    replay, replayCreated, err := service.Create(context.Background(),request)
    if err != nil || replayCreated || replay.GID!=createdPerson.GID { t.Fatalf("invalid idempotent replay: %v, %v",replayCreated,err) }
    if len(producer.changes)!=1 || work.runs!=1 { t.Fatal("existing Person must not enqueue another event") }
}

func TestPersonUpdateAndDeactivateEmitOnlyRealChanges(t *testing.T) {
    people, contacts, work, _ := setup(t)
    gid := "018f5d04-315c-75a9-b890-df0d02e787c1"
    people.person = entity.Person{GID:gid, FirstName:"Old", LastName:"Name", Status:entity.PersonStatusActive, Version:1, UpdatedAt:time.Now().UTC()}
    producer := &fakePersonChangeProducer{}
    service := personservice.NewPersonService(people, contacts, work, fakeMapper{}, fakeValidator{}, producer)
    name:="Updated"
    request:=persondto.UpdatePersonRequest{Params:persondto.PersonByGID{GID:gid},
        RequestBody:persondto.UpdatePerson{FirstName:annotationdto.Patch[string]{Set:true,Value:&name}}}
    if _,err:=service.Update(context.Background(),request);err!=nil{t.Fatal(err)}
    if len(producer.changes)!=1 || producer.changes[0].operation!=foundationevents.Updated ||
        producer.changes[0].old.Version!=1 || producer.changes[0].next.Version!=2 {
        t.Fatalf("incorrect Update event: %#v",producer.changes)
    }
    if _,err:=service.Update(context.Background(),request);err!=nil{t.Fatal(err)}
    if len(producer.changes)!=1 {t.Fatal("unchanged Update must not enqueue")}
    deactivation:=persondto.DeactivatePersonRequest{Params:persondto.PersonByGID{GID:gid}}
    if _,err:=service.Deactivate(context.Background(),deactivation);err!=nil{t.Fatal(err)}
    if len(producer.changes)!=2 || producer.changes[1].operation!=foundationevents.Deactivated {
        t.Fatalf("missing deactivation event: %#v",producer.changes)
    }
    if _,err:=service.Deactivate(context.Background(),deactivation);err!=nil{t.Fatal(err)}
    if len(producer.changes)!=2 || work.runs!=2 {t.Fatal("already deactivated must not enqueue or rewrite")}
}

func TestPersonEventEnqueueFailureFailsBusinessTransaction(t *testing.T) {
    people, contacts, work, _ := setup(t)
    producer:=&fakePersonChangeProducer{err:errors.New("outbox insert failed")}
    service:=personservice.NewPersonService(people,contacts,work,fakeMapper{},fakeValidator{},producer)
    _, created, err:=service.Create(context.Background(),persondto.CreatePersonRequest{RequestBody:persondto.CreatePerson{
        PersonNames:persondto.PersonNames{FirstName:"Anu",LastName:"Roy"},
    }})
    if err==nil || created || work.runs!=1 {t.Fatalf("failure must propagate from transaction: created=%v err=%v runs=%d",created,err,work.runs)}
}
