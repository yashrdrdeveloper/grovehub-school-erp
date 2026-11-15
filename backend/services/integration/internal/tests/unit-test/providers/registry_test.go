package providers_test

import (
    providerhelpers "grovehub/services/integration/internal/app/providers/helpers"
    register "grovehub/services/integration/internal/app/providers/operations/register"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/services/integration/internal/app/externalService/mail/gmail"
	"grovehub/services/integration/internal/app/externalService/mail/outlook"
)

var _ = Describe("Provider registry", func() {
	It("lists providers and resolves them case-insensitively", func() {
		registry, err := register.NewRegistry(gmail.New(), outlook.New())
		Expect(err).NotTo(HaveOccurred())
		Expect(registry.List()).To(HaveLen(2))
		provider, err := registry.Find("gmail")
		Expect(err).NotTo(HaveOccurred())
		Expect(provider.Descriptor().ID).To(Equal("GMAIL"))
	})
	It("rejects duplicate registrations", func() {
		_, err := register.NewRegistry(gmail.New(), gmail.New())
		Expect(err).To(MatchError(providerhelpers.ErrDuplicateProvider))
	})
	It("rejects unsupported providers", func() {
		registry, err := register.NewRegistry(gmail.New())
		Expect(err).NotTo(HaveOccurred())
		_, err = registry.Find("WHATSAPP")
		Expect(err).To(MatchError(providerhelpers.ErrUnknownProvider))
	})
})
