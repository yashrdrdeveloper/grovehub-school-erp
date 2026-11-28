package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/routes/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("DiscoverWithTier",func(){
	It("rejects an invalid handler signature",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"protected/record.go","package protected\nimport \"net/http\"\n// @Routes /records\ntype RecordController interface {\n// @GET /{gid}\nGet(http.ResponseWriter,http.Request)\n}\n")
		groups,err:=helpers.DiscoverWithTier(root,"/school",true,false)
		Expect(err).To(MatchError(ContainSubstring("must accept http.ResponseWriter and *http.Request")))
		Expect(groups).To(BeNil())
	})
})
