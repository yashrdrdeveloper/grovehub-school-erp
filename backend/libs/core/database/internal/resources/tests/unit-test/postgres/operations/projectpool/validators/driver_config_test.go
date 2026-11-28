package validators_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/postgres/operations/projectpool/validators"
)

var _=Describe("Driver configuration validation",func(){
	It("rejects invalid driver-only options and path identity mismatches",func(){
		for _,dsn:=range []string{
			"postgres://school@localhost/school?sslmode=invalid",
			"postgres://school@localhost/school?connect_timeout=invalid",
			"postgres://school@localhost//school",
		}{
			config,err:=validators.ProjectConfig(dsn)
			Expect(err).To(HaveOccurred(),dsn)
			Expect(config).To(BeNil())
		}
	})
})
