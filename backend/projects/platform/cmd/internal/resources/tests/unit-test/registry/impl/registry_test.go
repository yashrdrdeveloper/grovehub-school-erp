package impl_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"grovehub/cmd/internal/resources/tests/mocks/command"
	commandcontract "grovehub/cmd/internal/src/command"
	"grovehub/cmd/internal/src/registry/impl"
)

var (
	controller  *gomock.Controller
	mockCommand *command.MockCommand
	output      bytes.Buffer
	diagnostics bytes.Buffer
)

// @BeforeEach
func Reset() {
	controller = gomock.NewController(GinkgoT())
	mockCommand = command.NewMockCommand(controller)
	output.Reset()
	diagnostics.Reset()
}

// @Testing suite=PlatformCommandRegistry
func RegistrySpecs() {
	Describe("command dispatch", func() {
		It("defaults to the server strategy and forwards streams", func() {
			mockCommand.EXPECT().Run(gomock.Nil(), &output, &diagnostics).Return(7)
			registry := impl.New(map[string]commandcontract.Command{"serve": mockCommand})
			Expect(registry.Run(nil, &output, &diagnostics)).To(Equal(7))
		})

		It("forwards only command-specific arguments", func() {
			mockCommand.EXPECT().Run([]string{"all"}, &output, &diagnostics).Return(7)
			registry := impl.New(map[string]commandcontract.Command{"migrate": mockCommand})
			Expect(registry.Run([]string{"migrate", "all"}, &output, &diagnostics)).To(Equal(7))
		})

		It("rejects unknown and unconfigured commands", func() {
			registry := impl.New(map[string]commandcontract.Command{"nil": nil})
			for _, name := range []string{"unknown", "nil"} {
				diagnostics.Reset()
				Expect(registry.Run([]string{name}, &output, &diagnostics)).To(Equal(2))
				Expect(diagnostics.String()).To(ContainSubstring("usage:"))
			}
		})
	})
}
