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
}

func (stub *personServiceStub) Create(
	context.Context,
	foundationrequest.CreatePerson,
) (foundationresponse.Person, bool, error) {
	return stub.person, true, nil
}

func (stub *personServiceStub) GetByGID(
	context.Context,
	foundationrequest.PersonByGID,
) (foundationresponse.Person, error) {
	return foundationresponse.Person{}, nil
}

func (stub *personServiceStub) Update(
	context.Context,
	foundationrequest.PersonByGID,
	foundationrequest.UpdatePerson,
) (foundationresponse.Person, error) {
	return foundationresponse.Person{}, nil
}

func (stub *personServiceStub) Search(
	context.Context,
	coresearch.Request,
) (coresearch.Result[foundationresponse.Person], error) {
	return coresearch.Result[foundationresponse.Person]{}, nil
}

func (stub *personServiceStub) Deactivate(
	context.Context,
	foundationrequest.PersonByGID,
) (foundationresponse.Person, error) {
	return foundationresponse.Person{}, nil
}

var _ foundationservices.PersonService = (*personServiceStub)(nil)

var _ = Describe("Foundation person events", func() {
	It("publishes person created through infrastructure messaging", func() {
		infrastructureModule, err := infrastructure.New(infrastructureconfig.DefaultMessaging())
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			Expect(infrastructureModule.Close()).To(Succeed())
		})

		ctx, cancel := context.WithCancel(context.Background())
		DeferCleanup(cancel)

		gid := uuid.NewString()
		now := time.Now().UTC()
		person := foundationresponse.Person{
			GID:       gid,
			Status:    "ACTIVE",
			Version:   1,
			CreatedAt: now,
			UpdatedAt: now,
		}

		received := make(chan coreevents.Envelope[personevents.Created], 1)
		consumerDone := make(chan struct{})
		go func() {
			defer close(consumerDone)
			_ = infrastructureModule.MessagingClient().Consume(
				ctx,
				personevents.CreatedTopic,
				func(_ context.Context, message infrastructureservices.Message) error {
					var envelope coreevents.Envelope[personevents.Created]
					if err := json.Unmarshal(message.Payload, &envelope); err != nil {
						return nil
					}
					if envelope.Data.GID != gid {
						return nil
					}
					select {
					case received <- envelope:
					default:
					}
					return nil
				},
			)
		}()

		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		personService := personeventimpl.NewPersonService(
			&personServiceStub{person: person},
			infrastructureModule.MessagingClient(),
			logger,
		)

		createdPerson, created, err := personService.Create(ctx, foundationrequest.CreatePerson{})
		Expect(err).NotTo(HaveOccurred())
		Expect(created).To(BeTrue())
		Expect(createdPerson.GID).To(Equal(gid))

		var envelope coreevents.Envelope[personevents.Created]
		Eventually(received, sanityTimeout, sanityInterval).Should(Receive(&envelope))

		Expect(envelope.EventID).NotTo(BeEmpty())
		Expect(envelope.EventType).To(Equal(personevents.CreatedType))
		Expect(envelope.Source).To(Equal(personevents.Source))
		Expect(envelope.Data.GID).To(Equal(gid))
		Expect(envelope.Data.Status).To(Equal("ACTIVE"))
		Expect(envelope.Data.Version).To(Equal(int64(1)))

		cancel()
		Eventually(consumerDone, sanityTimeout, sanityInterval).Should(BeClosed())
	})
})
