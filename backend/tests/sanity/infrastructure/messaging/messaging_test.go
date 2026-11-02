package messaging_test

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/services/infrastructure"
	infrastructureconfig "grovehub/services/infrastructure/config"
	infrastructureservices "grovehub/services/infrastructure/services"
)

var _ = Describe("Infrastructure messaging", func() {
	var module infrastructure.Module

	BeforeEach(func() {
		var err error
		module, err = infrastructure.New(infrastructureconfig.DefaultMessaging())
		Expect(err).NotTo(HaveOccurred())

		DeferCleanup(func() {
			Expect(module.Close()).To(Succeed())
		})
	})

	It("reports NATS messaging as healthy", func() {
		Eventually(func() error {
			return module.MessagingClient().Health(context.Background())
		}, sanityTimeout, sanityInterval).Should(Succeed())
	})

	It("publishes and consumes a message through the infrastructure client", func() {
		ctx, cancel := context.WithCancel(context.Background())
		DeferCleanup(cancel)

		topic := "sanity.infrastructure.messaging." + uuid.NewString()
		payload := []byte("grovehub-messaging-sanity")
		metadata := map[string]string{"sanity": "true"}
		received := make(chan infrastructureservices.Message, 1)

		go func() {
			_ = module.MessagingClient().Consume(ctx, topic, func(_ context.Context, message infrastructureservices.Message) error {
				select {
				case received <- message:
				default:
				}
				return nil
			})
		}()

		Eventually(func() error {
			return module.MessagingClient().Publish(ctx, topic, payload, metadata)
		}, sanityTimeout, sanityInterval).Should(Succeed())

		var message infrastructureservices.Message
		Eventually(received, sanityTimeout, sanityInterval).Should(Receive(&message))

		Expect(message.ID).NotTo(BeEmpty())
		Expect(message.Payload).To(Equal(payload))
		Expect(message.Metadata).To(HaveKeyWithValue("sanity", "true"))
	})
})
