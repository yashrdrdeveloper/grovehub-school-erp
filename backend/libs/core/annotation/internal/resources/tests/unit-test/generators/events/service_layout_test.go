package events_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/events"
	sourceFixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Event output failures",func(){
	It("rejects stale nonempty destinations and looping source roots",func(){
		root:=GinkgoT().TempDir()
		owner:=filepath.Join(root,"services","demo")
		sourceFixture.Write(owner,"events/record/event.go","package record\n// @Event topic=record type=record.created.v1\ntype Created struct{}")
		Expect(events.GenerateAll(root)).To(Succeed())
		out:=filepath.Join(owner,".generated","events","record","zz_events.gen.go")
		Expect(os.Remove(out)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(out,"child"),0755)).To(Succeed())
		Expect(events.GenerateAll(root)).NotTo(Succeed())
		sourceFixture.Write(owner,"events/record/event.go","package record\ntype Created struct{}")
		Expect(events.GenerateAll(root)).NotTo(Succeed())
		Expect(os.RemoveAll(filepath.Join(owner,"events"))).To(Succeed())
		Expect(os.Symlink("events",filepath.Join(owner,"events"))).To(Succeed())
		Expect(events.GenerateAll(root)).NotTo(Succeed())
	})
})
