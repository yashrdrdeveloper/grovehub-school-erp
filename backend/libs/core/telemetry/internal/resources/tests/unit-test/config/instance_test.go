package config_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/telemetry/internal/src/config"
)

var _ = Describe("InstanceID", func() {
	It("uses a stable fallback for missing hostnames", func() {
		Expect(config.InstanceID(func() (string, error) { return "", errors.New("unavailable") }, 42)).To(Equal("unknown-host-42"))
		Expect(config.InstanceID(func() (string, error) { return " ", nil }, 42)).To(Equal("unknown-host-42"))
		Expect(config.InstanceID(func() (string, error) { return "host", nil }, 42)).To(Equal("host-42"))
	})
})
