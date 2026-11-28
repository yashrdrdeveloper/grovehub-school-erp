package cli_test

import (
	"bytes"
	"errors"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/cli"
)

type failedInventoryOutput struct{}
func (failedInventoryOutput) Write([]byte)(int,error){return 0,errors.New("output unavailable")}

var _=Describe("Inventory command errors",func(){
	It("reports missing owners and output failures",func(){
		var output,errOut bytes.Buffer
		root:=GinkgoT().TempDir()
		Expect(cli.Run([]string{"inventory",filepath.Join(root,"missing")},&output,&errOut)).To(Equal(1))
		Expect(errOut.String()).NotTo(BeEmpty())
		errOut.Reset()
		Expect(cli.Run([]string{"inventory",root},failedInventoryOutput{},&errOut)).To(Equal(1))
		Expect(errOut.String()).To(ContainSubstring("output unavailable"))
	})
})
