package providers_test

import (
    providerhelpers "grovehub/services/integration/internal/app/providers/helpers"
    mailconstants "grovehub/services/integration/internal/app/providers/mailprovider/constants"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	mailfactory "grovehub/services/integration/internal/app/providers/mailprovider/operations/build"
	"grovehub/services/integration/internal/app/externalService/mail/gmail"
	"grovehub/services/integration/internal/app/externalService/mail/outlook"
)

var _ = Describe("Mail provider factory", func() {
	It("resolves providers and their declared capabilities without hard-coded service dispatch", func() {
		runtime, err := mailfactory.NewFactory().Build(gmail.NewStrategy(nil), outlook.NewStrategy(nil))
		Expect(err).NotTo(HaveOccurred())
		Expect(runtime.Registry.List()).To(HaveLen(2))
		selection, err := runtime.Resolver.Resolve("gmail", mailconstants.CapabilityMailRead)
		Expect(err).NotTo(HaveOccurred())
		Expect(selection.Descriptor().ID).To(Equal("GMAIL"))
		_, err = runtime.Resolver.Resolve("OUTLOOK", mailconstants.CapabilityMailWatch)
		Expect(err).To(MatchError(providerhelpers.ErrUnsupportedAction))
		Expect(runtime.Actions.Actions("GMAIL")).To(BeEmpty()) // no OAuth configuration
	})
	It("rejects duplicate provider strategies", func() {
		_, err := mailfactory.NewFactory().Build(gmail.NewStrategy(nil), gmail.NewStrategy(nil))
		Expect(err).To(MatchError(providerhelpers.ErrDuplicateProvider))
	})
})
