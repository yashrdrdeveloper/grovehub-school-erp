package events_test

import (
    "context"
    "encoding/json"
    "io"
    "log/slog"
    "time"

    "github.com/google/uuid"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    coreevents "grovehub/libs/core/events"
    coresearch "grovehub/libs/core/search"
    foundationevents "grovehub/services/foundation/constants/events"
    foundationrequest "grovehub/services/foundation/dto/request"
    foundationresponse "grovehub/services/foundation/dto/response"
    personevents "grovehub/services/foundation/events/persons"
    personeventimpl "grovehub/services/foundation/events/persons/impl"
    foundationservices "grovehub/services/foundation/services"
    "grovehub/services/infrastructure"
    infrastructureconfig "grovehub/services/infrastructure/config"
    infrastructureservices "grovehub/services/infrastructure/services"
)

type personServiceStub struct {
    person foundationresponse.Person
    original foundationresponse.Person
    updated foundationresponse.Person
    created bool
}
func (stub *personServiceStub) Create(context.Context, foundationrequest.CreatePerson) (foundationresponse.Person, bool, error) {
    return stub.person, stub.created, nil
}
func (stub *personServiceStub) GetByGID(context.Context, foundationrequest.PersonByGID) (foundationresponse.Person, error) {
    return stub.original, nil
}
func (stub *personServiceStub) Update(context.Context, foundationrequest.PersonByGID, foundationrequest.UpdatePerson) (foundationresponse.Person, error) {
    return stub.updated, nil
}
func (stub *personServiceStub) Search(context.Context, coresearch.Request) (coresearch.Result[foundationresponse.Person], error) {
    return coresearch.Result[foundationresponse.Person]{}, nil
}
func (stub *personServiceStub) Deactivate(context.Context, foundationrequest.PersonByGID) (foundationresponse.Person, error) {
    return stub.updated, nil
}
var _ foundationservices.PersonService = (*personServiceStub)(nil)

var _ = Describe("Foundation person changed events", func() {
    It("publishes CREATED, UPDATED and DEACTIVATED snapshots using one topic", func() {
        infrastructureModule, err := infrastructure.New(infrastructureconfig.DefaultMessaging())
        Expect(err).NotTo(HaveOccurred())
        DeferCleanup(func() { Expect(infrastructureModule.Close()).To(Succeed()) })

        ctx, cancel := context.WithCancel(context.Background())
        DeferCleanup(cancel)
        gid := uuid.NewString()
        now := time.Now().UTC()
        old := foundationresponse.Person{GID:gid, FirstName:"Arun", LastName:"Kumar", Status:"ACTIVE", Version:1, CreatedAt:now, UpdatedAt:now}
        next := old
        next.LastName = "Sharma"
        next.Version = 2
        deactivated := next
        deactivated.Status = "DEACTIVATED"
        deactivated.Version = 3

        received := make(chan coreevents.Envelope[personevents.Changed], 3)
        consumerDone := make(chan struct{})
        go func() {
            defer close(consumerDone)
            _ = infrastructureModule.MessagingClient().Consume(ctx, personevents.ChangedTopic,
                func(_ context.Context, message infrastructureservices.Message) error {
                    var envelope coreevents.Envelope[personevents.Changed]
                    if err := json.Unmarshal(message.Payload, &envelope); err != nil { return err }
                    if envelope.Data.New == nil || envelope.Data.New.GID != gid { return nil }
                    select { case received <- envelope: default: }
                    return nil
                })
        }()

        logger := slog.New(slog.NewTextHandler(io.Discard, nil))
        wrap := func(stub *personServiceStub) foundationservices.PersonService {
            return personeventimpl.NewPersonService(stub, infrastructureModule.MessagingClient(), logger)
        }
        createdPerson, created, err := wrap(&personServiceStub{person:old, created:true}).Create(ctx, foundationrequest.CreatePerson{})
        Expect(err).NotTo(HaveOccurred())
        Expect(created).To(BeTrue())
        Expect(createdPerson.GID).To(Equal(gid))

        changedPerson, err := wrap(&personServiceStub{original:old, updated:next}).Update(ctx, foundationrequest.PersonByGID{GID:gid}, foundationrequest.UpdatePerson{})
        Expect(err).NotTo(HaveOccurred())
        Expect(changedPerson.Version).To(Equal(int64(2)))

        changedPerson, err = wrap(&personServiceStub{original:next, updated:deactivated}).Deactivate(ctx, foundationrequest.PersonByGID{GID:gid})
        Expect(err).NotTo(HaveOccurred())
        Expect(changedPerson.Status).To(Equal("DEACTIVATED"))

        got := map[foundationevents.Operation]coreevents.Envelope[personevents.Changed]{}
        Eventually(func() int {
            for {
                select {
                case envelope := <-received: got[envelope.Data.Operation] = envelope
                default: return len(got)
                }
            }
        }, sanityTimeout, sanityInterval).Should(Equal(3))
        for _, op := range []foundationevents.Operation{foundationevents.Created, foundationevents.Updated, foundationevents.Deactivated} {
            envelope := got[op]
            Expect(envelope.EventID).NotTo(BeEmpty())
            Expect(envelope.EventType).To(Equal(personevents.ChangedType))
            Expect(envelope.Version).To(Equal(1))
            Expect(envelope.AggregateID).To(Equal(gid))
            Expect(envelope.Data.New.GID).To(Equal(gid))
        }
        Expect(got[foundationevents.Created].Data.Old).To(BeNil())
        Expect(got[foundationevents.Updated].Data.Old.LastName).To(Equal("Kumar"))
        Expect(got[foundationevents.Updated].Data.New.LastName).To(Equal("Sharma"))
        Expect(got[foundationevents.Deactivated].Data.Old.Status).To(Equal("ACTIVE"))
        Expect(got[foundationevents.Deactivated].Data.New.Status).To(Equal("DEACTIVATED"))

        cancel()
        Eventually(consumerDone, sanityTimeout, sanityInterval).Should(BeClosed())
    })
})
