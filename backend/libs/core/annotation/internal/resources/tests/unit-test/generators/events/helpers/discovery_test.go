package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/events/helpers"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Discover", func() {
	It("rejects conflicting source declarations within one event package",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"a.go","package events\n// @Event topic=person type=person.created.v1 source=school\ntype Created struct{}\n")
		fixture.Write(root,"b.go","package events\n// @Event topic=person type=person.changed.v1 source=platform\ntype Changed struct{}\n")
		_,_,err:=helpers.Discover(root)
		Expect(err).To(MatchError(ContainSubstring("conflicting event sources")))
	})
})
