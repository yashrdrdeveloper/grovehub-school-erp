package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/events/helpers"
)

var _ = Describe("Event comment discovery", func() {
	It("ignores descriptive comments", func() {
		root := GinkgoT().TempDir()
		fixture.Write(root, "event.go", "package events\n// Created describes an event.\n// @Event topic=person type=person.created.v1 source=school\ntype Created struct{}")
		defs, _, err := helpers.Discover(root)
		Expect(err).NotTo(HaveOccurred())
		Expect(defs).To(HaveLen(1))
	})
})
