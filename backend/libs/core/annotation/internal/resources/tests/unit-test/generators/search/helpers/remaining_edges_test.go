package helpers_test

import (
"os"
"path/filepath"
"strings"
"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/search/helpers"
fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"
)

var _ = Describe("Search input boundary failures",func(){
 It("reports a blocked implementation directory",func(){
 root:=searchFixture(searchTestImpl,searchTestConfig)
 Expect(os.RemoveAll(filepath.Join(root,"impl"))).To(Succeed())
 fixture.Write(root,"impl","blocked")
 _,err:=helpers.Discover(root,"grovehub","demo");Expect(err).To(HaveOccurred())
 })
 It("reports a blocked configuration package directory",func(){
 root:=searchFixture(searchTestImpl,searchTestConfig)
 fixture.Write(root,"search","blocked")
 _,err:=helpers.Discover(root,"grovehub","demo");Expect(err).To(HaveOccurred())
 })
 It("skips unrelated declarations when resolving config",func(){
 root:=searchFixture(searchTestImpl,searchTestConfig+"\nvar X,Y int\n")
 _,err:=helpers.Discover(root,"grovehub","demo");Expect(err).NotTo(HaveOccurred())
 })
 It("rejects a missing static field",func(){
 root:=searchFixture(searchTestImpl,strings.Replace(searchTestConfig,"CursorExpression:","Unrelated:",1))
 _,err:=helpers.Discover(root,"grovehub","demo");Expect(err).To(HaveOccurred())
 })
 It("rejects an interface without its repository directive",func(){
 root:=searchFixture(searchTestImpl,searchTestConfig)
 fixture.Write(root,"person.go",strings.Replace(searchTestInterface,"// @Repository entity=Person\n","",1))
 _,err:=helpers.Discover(root,"grovehub","demo");Expect(err).To(HaveOccurred())
 })
 It("rejects a concrete declaration without a body",func(){
 root:=searchFixture(strings.Split(searchTestImpl," {\n\treturn repository.searchGenerated")[0],searchTestConfig)
 _,err:=helpers.Discover(root,"grovehub","demo");Expect(err).To(HaveOccurred())
 })
 It("skips aliases while validating implementation structs",func(){
 root:=searchFixture(searchTestImpl+"\ntype Alias int\n",searchTestConfig)
 specs,err:=helpers.Discover(root,"grovehub","demo");Expect(err).NotTo(HaveOccurred())
 Expect(helpers.ValidateImplementation(filepath.Join(root,"impl"),specs)).To(Succeed())
 })
})
