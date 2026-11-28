package configuration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
	"os"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/configuration"
)

type settings struct {
	Host string `env:"HOST" required:"true"`
	Port int `env:"PORT" default:"8080"`
	Enabled bool `env:"ENABLED"`
	Timeout time.Duration `env:"TIMEOUT"`
}

var _ = Describe("Decode", func() {
	It("decodes typed environment values and defaults", func() {
		values := map[string]string{"HOST":" localhost ", "ENABLED":"true", "TIMEOUT":"3s"}
		value, err := configuration.Decode[settings](func(key string) (string, bool) { v, ok := values[key]; return v, ok })
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal(settings{Host:"localhost", Port:8080, Enabled:true, Timeout:3*time.Second}))
	})
	It("rejects missing required fields and invalid targets", func() {
		_, err := configuration.Decode[settings](func(string)(string,bool){return "",false})
		Expect(err).To(MatchError(ContainSubstring("required configuration HOST")))
		_, err = configuration.Decode[settings](nil)
		Expect(err).To(HaveOccurred())
		_, err = configuration.Decode[string](os.LookupEnv)
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("FromEnv", func() {
	It("reads the actual environment", func() {
		GinkgoT().Setenv("ANNOTATION_TEST_HOST", "localhost")
		type input struct { Host string `env:"ANNOTATION_TEST_HOST"` }
		value, err := configuration.FromEnv[input]()
		Expect(err).NotTo(HaveOccurred())
		Expect(value.Host).To(Equal("localhost"))
	})
})
