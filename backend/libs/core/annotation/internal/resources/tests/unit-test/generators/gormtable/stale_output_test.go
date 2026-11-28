package gormtable_test

import (
	"os"
	"path/filepath"
	

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/gormtable"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Table source and stale output handling",func(){
	It("ignores unrelated function and documentation declarations and rejects nonempty stale output directories",func(){
		root,out:=GinkgoT().TempDir(),GinkgoT().TempDir()
		fixture.Write(root,"record.go","package entity\nfunc Helper(){}\n// explanatory text\ntype Plain struct{}")
		path:=filepath.Join(out,"zz_tables.gen.go")
		Expect(os.MkdirAll(filepath.Join(path,"child"),0755)).To(Succeed())
		Expect(gormtable.GenerateInto(root,path)).NotTo(Succeed())
		Expect(path).To(BeADirectory())
	})
})
