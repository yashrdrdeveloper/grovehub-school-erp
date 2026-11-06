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
	"grovehub/services/identity/entity"
	accountevents "grovehub/services/identity/events/accounts"
	accounteventimpl "grovehub/services/identity/events/accounts/impl"
	identityservices "grovehub/services/identity/services"
	"grovehub/services/infrastructure"
	infrastructureconfig "grovehub/services/infrastructure/config"
	infrastructureservices "grovehub/services/infrastructure/services"
)

type accountServiceStub struct {
	account entity.Account
}

func (stub *accountServiceStub) Create(
	context.Context,
	identityservices.CreateAccountInput,
) (entity.Account, bool, error) {
	return stub.account, true, nil
}

func (stub *accountServiceStub) GetById(context.Context, string) (entity.Account, error) {
	return entity.Account{}, nil
}

func (stub *accountServiceStub) GetByPersonGid(context.Context, string) (entity.Account, error) {
	return entity.Account{}, nil
}

func (stub *accountServiceStub) GetByEmail(context.Context, string) (entity.Account, error) {
	return entity.Account{}, nil
}

func (stub *accountServiceStub) GetByRegistrationIdempotencyKey(context.Context, string) (entity.Account, error) {
	return entity.Account{}, nil
}

func (stub *accountServiceStub) Update(
	context.Context,
	string,
	identityservices.UpdateAccountInput,
) (entity.Account, error) {
	return entity.Account{}, nil
}

func (stub *accountServiceStub) Deactivate(context.Context, string) (entity.Account, error) {
	return entity.Account{}, nil
}

func (stub *accountServiceStub) PermanentlyDelete(context.Context, string) error {
	return nil
}

var _ = Describe("Identity account events", func() {
	It("publishes account created through infrastructure messaging", func() {
		infrastructureModule, err := infrastructure.New(infrastructureconfig.DefaultMessaging())
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			Expect(infrastructureModule.Close()).To(Succeed())
		})

		ctx, cancel := context.WithCancel(context.Background())
		DeferCleanup(cancel)

		accountID := uuid.NewString()
		personGID := uuid.NewString()
		now := time.Now().UTC()
		account := entity.NewAccount(
			entity.AccountId(accountID),
			entity.PersonGid(personGID),
			"sanity@example.invalid",
			"",
			entity.AccountStatusActive,
			1,
			now,
			now,
		)

		received := make(chan coreevents.Envelope[accountevents.Created], 1)
		consumerDone := make(chan struct{})
		go func() {
			defer close(consumerDone)
			_ = infrastructureModule.MessagingClient().Consume(
				ctx,
				accountevents.CreatedTopic,
				func(_ context.Context, message infrastructureservices.Message) error {
					var envelope coreevents.Envelope[accountevents.Created]
					if err := json.Unmarshal(message.Payload, &envelope); err != nil {
						return nil
					}
					if envelope.Data.AccountID != accountID {
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
		accountService := accounteventimpl.NewAccountService(
			&accountServiceStub{account: account},
			infrastructureModule.MessagingClient(),
			logger,
		)

		createdAccount, created, err := accountService.Create(ctx, identityservices.CreateAccountInput{})
		Expect(err).NotTo(HaveOccurred())
		Expect(created).To(BeTrue())
		Expect(string(createdAccount.GetAccountId())).To(Equal(accountID))

		var envelope coreevents.Envelope[accountevents.Created]
		Eventually(received, sanityTimeout, sanityInterval).Should(Receive(&envelope))

		Expect(envelope.EventID).NotTo(BeEmpty())
		Expect(envelope.EventType).To(Equal(accountevents.CreatedType))
		Expect(envelope.Version).To(Equal(1))
		Expect(envelope.AggregateID).To(Equal(accountID))
		Expect(envelope.Data.AccountID).To(Equal(accountID))
		Expect(envelope.Data.PersonGID).To(Equal(personGID))
		Expect(envelope.Data.Status).To(Equal(string(entity.AccountStatusActive)))
		Expect(envelope.Data.Version).To(Equal(int64(1)))

		cancel()
		Eventually(consumerDone, sanityTimeout, sanityInterval).Should(BeClosed())
	})
})
