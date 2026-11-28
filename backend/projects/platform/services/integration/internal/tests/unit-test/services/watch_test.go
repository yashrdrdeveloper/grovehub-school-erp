package services_test

import (
    register "grovehub/services/integration/internal/app/providers/operations/register"
    execute "grovehub/services/integration/internal/app/providers/operations/execute"
    oauth "grovehub/services/integration/internal/app/providers/mailprovider/operations/oauth"
	"context"
    schoolgateway "grovehub/libs/platform/session"
	"encoding/json"
    "errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	messagingclient "grovehub/services/infrastructure/client/messaging"
 integrationconfig "grovehub/services/integration/internal/app/config"
	"grovehub/services/integration/internal/app/entity"
	mailboxevents "grovehub/services/integration/internal/app/events/producers/connection/operations/mailbox"
	
	"grovehub/services/integration/internal/app/externalService/mail/gmail"
	integrationimpl "grovehub/services/integration/internal/app/services/connection/impl"
    changed "grovehub/services/integration/internal/app/events/producers/connection/operations/mailbox"
    changedimpl "grovehub/services/integration/internal/app/events/producers/connection/operations/mailbox/strategies/impl"
    producerimpl "grovehub/services/integration/internal/app/events/producers/connection/impl"
    validatorimpl "grovehub/services/integration/internal/app/validators/connection/impl"
    mapperimpl "grovehub/services/integration/internal/app/mappers/connection/impl"
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

func (*connectionRepositoryFake) SystemMailSender(context.Context, string) (entity.Connection, error) {
	return entity.Connection{}, nil
}

func (*connectionRepositoryFake) ConfigureSystemMailSender(context.Context, string, string) error {
	return nil
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

type messagingFake struct{
    published []publishedEvent
    failure error
}

func (m *messagingFake) Publish(_ context.Context, topic string, payload []byte, _ map[string]string) error {
    if m.failure!=nil {return m.failure}
	m.published = append(m.published, publishedEvent{topic: topic, data: payload})
	return nil
}

func (m *messagingFake) PublishWithOptions(ctx context.Context, topic string, payload []byte, metadata map[string]string, options messagingclient.PublishOptions) error {
    return m.Publish(ctx,topic,payload,metadata)
}

func (*messagingFake) Consume(context.Context, string, messagingclient.MessageHandler) error {
	return nil
}

func (*messagingFake) Health(context.Context) error { return nil }

var _ = Describe("Provider push event ownership", func() {
	It("publishes one normalized event per connection without storing or syncing messages", func() {
		registry, err := register.NewRegistry(gmail.New())
		Expect(err).NotTo(HaveOccurred())
		messaging := &messagingFake{}
        producer:=producerimpl.New(changed.NewChangedPipeline(changedimpl.New(),messaging))
		service, err := integrationimpl.NewConnectionService(
			&connectionRepositoryFake{}, &watchRepositoryFake{},
			registry, execute.NewActionRegistry(), oauth.NewRegistry(), gmail.NewClient(integrationconfig.Gmail{}), nil, []byte("01234567890123456789012345678901"), "projects/test/topics/gmail", producer,validatorimpl.NewConnectionValidator(),mapperimpl.NewConnectionMapper(),
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
    It("returns NATS delivery errors so the Gmail push endpoint can retry",func(){
        registry,err:=register.NewRegistry(gmail.New())
        Expect(err).NotTo(HaveOccurred())
        sentinel:=errors.New("broker unavailable")
        transport:=&messagingFake{failure:sentinel}
        producer:=producerimpl.New(changed.NewChangedPipeline(changedimpl.New(),transport))
        service,err:=integrationimpl.NewConnectionService(
            &connectionRepositoryFake{},&watchRepositoryFake{},
            registry,execute.NewActionRegistry(),oauth.NewRegistry(),
            gmail.NewClient(integrationconfig.Gmail{}),nil,
            []byte("01234567890123456789012345678901"),
            "projects/test/topics/gmail",producer,
            validatorimpl.NewConnectionValidator(),mapperimpl.NewConnectionMapper(),
        )
        Expect(err).NotTo(HaveOccurred())
        err=service.ProcessPush(context.Background(),gmail.PushNotification{
            EmailAddress:"one@gmail.com",HistoryID:"10001",
        })
        Expect(err).To(MatchError(sentinel))
        Expect(transport.published).To(BeEmpty())
    })

})

var _ = Describe("School lifecycle for provider push",func(){
 It("skips disabled tenants and retries lookup failures without blocking active tenants",func(){
  registry,err:=register.NewRegistry(gmail.New())
  Expect(err).NotTo(HaveOccurred())
  transport:=&messagingFake{}
  producer:=producerimpl.New(changed.NewChangedPipeline(changedimpl.New(),transport))
  service,err:=integrationimpl.NewConnectionService(
   &connectionRepositoryFake{},&watchRepositoryFake{},
   registry,execute.NewActionRegistry(),oauth.NewRegistry(),
   gmail.NewClient(integrationconfig.Gmail{}),nil,
   []byte("01234567890123456789012345678901"),
   "projects/test/topics/gmail",producer,
   validatorimpl.NewConnectionValidator(),mapperimpl.NewConnectionMapper(),
  )
  Expect(err).NotTo(HaveOccurred())
  unavailable:=false
  sentinel:=errors.New("registry unavailable")
  guarded:=integrationimpl.WithTenantState(service,func(_ context.Context,id string)(schoolgateway.TenantState,error){
   if id=="tenant-1"{
    if unavailable{return schoolgateway.TenantState{},sentinel}
    return schoolgateway.TenantState{Generation:2,Active:false},nil
   }
   return schoolgateway.TenantState{Generation:1,Active:true},nil
  })
  push:=gmail.PushNotification{EmailAddress:"one@gmail.com",HistoryID:"10001"}
  Expect(guarded.ProcessPush(context.Background(),push)).To(Succeed())
  Expect(transport.published).To(HaveLen(1))
  var envelope mailboxevents.Changed
  Expect(json.Unmarshal(transport.published[0].data,&envelope)).To(Succeed())
  Expect(envelope.Data.ConnectionID).To(Equal("connection-2"))
  unavailable=true
  Expect(guarded.ProcessPush(context.Background(),push)).To(MatchError(ContainSubstring("registry unavailable")))
  Expect(transport.published).To(HaveLen(2))
  Expect(json.Unmarshal(transport.published[1].data,&envelope)).To(Succeed())
  Expect(envelope.Data.ConnectionID).To(Equal("connection-2"))
 })
})
