package services_test

import (
	"context"
	"encoding/json"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	messagingclient "grovehub/services/infrastructure/client"
	"grovehub/services/integration/entity"
	mailboxevents "grovehub/services/integration/events/connections/mailbox"
	providerimpl "grovehub/services/integration/providers/impl"
	"grovehub/services/integration/providers/mailprovider/providers/gmail"
	integrationimpl "grovehub/services/integration/services/impl"
)

type connectionRepositoryFake struct{}

func (*connectionRepositoryFake) Create(context.Context, entity.Connection) error { return nil }
func (*connectionRepositoryFake) Get(context.Context, string, string) (entity.Connection, error) {
	return entity.Connection{}, nil
}
func (*connectionRepositoryFake) GetByID(context.Context, string) (entity.Connection, error) {
	return entity.Connection{}, nil
}
func (*connectionRepositoryFake) List(context.Context, string, string) ([]entity.Connection, error) {
	return nil, nil
}
func (*connectionRepositoryFake) Disconnect(context.Context, string, string) (entity.Connection, error) {
	return entity.Connection{}, nil
}
func (*connectionRepositoryFake) SaveOAuthState(context.Context, entity.OAuthState) error { return nil }
func (*connectionRepositoryFake) ConsumeOAuthState(context.Context, string, time.Time) (entity.OAuthState, error) {
	return entity.OAuthState{}, nil
}
func (*connectionRepositoryFake) Activate(context.Context, string, string, string, []byte) error {
	return nil
}
func (*connectionRepositoryFake) Credential(context.Context, string) ([]byte, error)     { return nil, nil }
func (*connectionRepositoryFake) UpdateCredential(context.Context, string, []byte) error { return nil }
func (*connectionRepositoryFake) SaveTest(context.Context, entity.ConnectionTest) error  { return nil }
func (*connectionRepositoryFake) ListTests(context.Context, string) ([]entity.ConnectionTest, error) {
	return nil, nil
}

type watchRepositoryFake struct{}

func (*watchRepositoryFake) SaveWatch(context.Context, string, string, time.Time) error { return nil }
func (*watchRepositoryFake) GetGmailByMailbox(_ context.Context, email string) ([]entity.Connection, error) {
	if email != "one@gmail.com" {
		return nil, nil
	}
	return []entity.Connection{
		{ID: "connection-1", TenantID: "tenant-1", OwnerGID: "person-1", Provider: "GMAIL", Status: entity.StatusActive},
		{ID: "connection-2", TenantID: "tenant-2", OwnerGID: "person-2", Provider: "GMAIL", Status: entity.StatusActive},
	}, nil
}
func (*watchRepositoryFake) ExpiringWatches(context.Context, time.Time) ([]entity.Connection, error) {
	return nil, nil
}

type publishedEvent struct {
	topic string
	data  []byte
}
type messagingFake struct{ published []publishedEvent }

func (m *messagingFake) Publish(_ context.Context, topic string, payload []byte, _ map[string]string) error {
	m.published = append(m.published, publishedEvent{topic: topic, data: payload})
	return nil
}
func (*messagingFake) Consume(context.Context, string, messagingclient.MessageHandler) error {
	return nil
}
func (*messagingFake) Health(context.Context) error { return nil }

var _ = Describe("Provider push event ownership", func() {
	It("publishes one normalized event per connection without storing or syncing messages", func() {
		registry, err := providerimpl.NewRegistry(gmail.New())
		Expect(err).NotTo(HaveOccurred())
		messaging := &messagingFake{}
		service, err := integrationimpl.NewConnectionService(
			&connectionRepositoryFake{}, &watchRepositoryFake{},
			registry, providerimpl.NewActionRegistry(), providerimpl.NewOAuthRegistry(), nil, nil, nil, "projects/test/topics/gmail", messaging,
		)
		Expect(err).NotTo(HaveOccurred())
		err = service.ProcessPush(context.Background(), gmail.PushNotification{
			EmailAddress: "one@gmail.com", HistoryID: "10001",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(messaging.published).To(HaveLen(2))
		ids := map[string]bool{}
		for _, published := range messaging.published {
			Expect(published.topic).To(Equal(mailboxevents.ChangedTopic))
			var envelope mailboxevents.Changed
			Expect(json.Unmarshal(published.data, &envelope)).To(Succeed())
			Expect(envelope.EventType).To(Equal(mailboxevents.ChangedEventType))
			Expect(envelope.Data.HistoryHint).To(Equal("10001"))
			ids[envelope.Data.ConnectionID] = true
		}
		Expect(ids).To(HaveKey("connection-1"))
		Expect(ids).To(HaveKey("connection-2"))
	})
})
