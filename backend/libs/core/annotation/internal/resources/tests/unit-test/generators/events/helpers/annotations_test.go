package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/events/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Discover", func() {
	DescribeTable("rejects malformed event metadata",func(args string){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"event.go","package events\n// @Event "+args+"\ntype Event struct{}\n")
		_,_,err:=helpers.Discover(root)
		Expect(err).To(HaveOccurred())
	},Entry("version zero","topic=person type=person.created.v0"),Entry("uppercase topic","topic=Person type=person.created.v1"),Entry("unknown key","topic=person type=person.created.v1 extra=x"),Entry("duplicate key","topic=person topic=other type=person.created.v1"),Entry("invalid source","topic=person type=person.created.v1 source=BAD"))
})
