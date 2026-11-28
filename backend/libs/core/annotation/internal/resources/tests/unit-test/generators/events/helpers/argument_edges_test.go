package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/events/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Event argument validation",func(){
	It("rejects arguments without an equals sign",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"event.go","package event\n// @Event topic=record missing type=record.created.v1\ntype Created struct{}")
		_,_,err:=helpers.Discover(root)
		Expect(err).To(MatchError(ContainSubstring("invalid event argument")))
	})
	It("permits events without a source alongside a declared source",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"event.go","package event\n// @Event topic=record type=record.created.v1\ntype Created struct{}\n// @Event topic=record type=record.updated.v2 source=school\ntype Updated struct{}")
		defs,pkg,err:=helpers.Discover(root)
		Expect(err).NotTo(HaveOccurred());Expect(defs).To(HaveLen(2));Expect(pkg).To(Equal("event"))
	})
})
