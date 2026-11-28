package validators_test

import (
	"github.com/jackc/pgx/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/postgres/operations/projectpool/validators"
)

var _=Describe("Driver configuration validation",func(){
	It("rejects invalid driver-only options and path identity mismatches",func(){
		for _,dsn:=range []string{
			"postgres://school@localhost/school?sslmode=invalid",
			"postgres://school@localhost/school?connect_timeout=invalid",
		}{
			config,err:=validators.ProjectConfig(dsn)
			Expect(err).To(HaveOccurred(),dsn)
			Expect(config).To(BeNil())
		}
	})
	It("rejects a driver-resolved identity that differs from the explicit URI",func(){
		const dsn="postgres://school@localhost/school?sslmode=disable"
		for _,change:=range []func(*pgx.ConnConfig){
			func(config *pgx.ConnConfig){config.User="other"},
			func(config *pgx.ConnConfig){config.Database="other"},
		}{
			config,err:=validators.ProjectConfigWithParser(dsn,func(raw string)(*pgx.ConnConfig,error){
				resolved,err:=pgx.ParseConfig(raw)
				if err!=nil{return nil,err}
				change(resolved)
				return resolved,nil
			})
			Expect(err).To(HaveOccurred())
			Expect(config).To(BeNil())
		}
	})
})
