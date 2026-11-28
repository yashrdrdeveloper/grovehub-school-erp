package events_test

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/events"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("GenerateAll",func(){
	It("generates event constants and removes stale outputs when annotations disappear",func(){
		root:=GinkgoT().TempDir()
		path:="services/demo/internal/app/events/record/event.go"
		fixture.Write(root,path,"package record\n// @Event topic=person type=person.created.v1\ntype Created struct{}\n")
		Expect(events.GenerateAll(root)).To(Succeed())
		out:=filepath.Join(root,"services/demo/.generated/internal/app/events/record/zz_events.gen.go")
		Expect(fixture.Read(out)).To(ContainSubstring("CreatedVersion"))
		fixture.Write(root,path,"package record\ntype Created struct{}\n")
		Expect(events.GenerateAll(root)).To(Succeed())
		Expect(out).NotTo(BeAnExistingFile())
	})
})
