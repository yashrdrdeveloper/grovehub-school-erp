package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/events/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Event declaration validation",func(){
	It("rejects duplicate declarations, unsupported options and invalid identifiers",func(){
		for _,test:=range []struct{annotation,reason string}{
			{"@Event topic=record type=record.created.v1\n// @Event topic=record type=record.created.v1","duplicate @Event"},
			{"@Event topic=record topic=other type=record.created.v1","invalid event argument"},
			{"@Event topic= type=record.created.v1","invalid event argument"},
			{"@Event topic=record type=record.created.v1 extra=unknown","unknown @Event key"},
			{"@Event topic=Record type=record.created.v1","requires valid topic"},
			{"@Event topic=record type=record.created.v0","requires valid topic"},
			{"@Event topic=record type=record.created.v1 source=Bad","invalid source"},
			{"@Event topic=record type=record.created.v1 typeConst=bad-name","invalid type constant"},
			{"@Event topic=record type=record.created.v99999999999999999999999999","value out of range"},
		}{
			root:=GinkgoT().TempDir()
			fixture.Write(root,"event.go","package event\n// "+test.annotation+"\ntype Created struct{}")
			_,_,err:=helpers.Discover(root)
			Expect(err).To(MatchError(ContainSubstring(test.reason)),test.annotation)
		}
	})
	It("rejects mixed packages and conflicting source ownership and ignores test files",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"ignored_test.go","invalid source")
		fixture.Write(root,"event.go","package event\n// @Event topic=record type=record.created.v1 source=school\ntype Created struct{}\n// @Event topic=record type=record.updated.v2 source=other\ntype Updated struct{}")
		_,_,err:=helpers.Discover(root)
		Expect(err).To(MatchError(ContainSubstring("conflicting event sources")))
		fixture.Write(root,"event.go","package event\nconst x=1\ntype Plain struct{}")
		fixture.Write(root,"other.go","package other")
		_,_,err=helpers.Discover(root)
		Expect(err).To(MatchError(ContainSubstring("mixed packages")))
		fixture.Write(root,"other.go","package event\nfunc broken(")
		_,_,err=helpers.Discover(root)
		Expect(err).To(HaveOccurred())
	})
})
