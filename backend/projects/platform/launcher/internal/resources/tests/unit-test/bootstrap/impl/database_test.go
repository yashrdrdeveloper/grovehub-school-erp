package impl_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
	"grovehub/launcher/internal/src/bootstrap/impl"
)

// @Testing suite=LauncherBootstrapImplementation
func DatabaseSpecs(){
	Describe("PostgreSQL pool ownership",func(){
		It("rejects missing DSN, bad pool limits, failed GORM open, and a missing underlying SQL pool",func(){
			t:=GinkgoT()
			ctx:=context.Background()
			t.Setenv("DATABASE_URL","")
			t.Setenv("PLATFORM_DB_MAX_OPEN","")
			t.Setenv("PLATFORM_DB_MAX_IDLE","")
			pools:=impl.NewPostgresPools(func(string)(*gorm.DB,error){Fail("unexpected opener");return nil,nil})
			_,err:=pools.Business(ctx)
			Expect(err).To(MatchError(impl.ErrDatabaseURLRequired))
			t.Setenv("DATABASE_URL","fixture")
			t.Setenv("PLATFORM_DB_MAX_OPEN","0")
			_,err=pools.Business(ctx)
			Expect(err).To(HaveOccurred())
			t.Setenv("PLATFORM_DB_MAX_OPEN","")
			pools=impl.NewPostgresPools(func(string)(*gorm.DB,error){return nil,errors.New("open failed")})
			_,err=pools.Business(ctx)
			Expect(err).To(MatchError(ContainSubstring("create postgres gorm db")))
			pools=impl.NewPostgresPools(func(string)(*gorm.DB,error){return &gorm.DB{Config:&gorm.Config{}},nil})
			_,err=pools.Business(ctx)
			Expect(err).To(MatchError(ContainSubstring("access postgres sql db")))
		})
		It("configures a healthy business pool and closes it after an unsuccessful health check",func(){
			t:=GinkgoT()
			t.Setenv("DATABASE_URL"," fixture ")
			t.Setenv("PLATFORM_DB_MAX_OPEN","7")
			t.Setenv("PLATFORM_DB_MAX_IDLE","3")
			ctx:=context.Background()
			db,expect:=mockPool()
			expect.ExpectPing()
			pools:=impl.NewPostgresPools(func(dsn string)(*gorm.DB,error){Expect(dsn).To(Equal("fixture"));return db,nil})
			got,err:=pools.Business(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(BeIdenticalTo(db))
			raw,err:=got.DB()
			Expect(err).NotTo(HaveOccurred())
			Expect(raw.Stats().MaxOpenConnections).To(Equal(7))
			Expect(expect.ExpectationsWereMet()).To(Succeed())
			failed,failedExpect:=mockPool()
			failedExpect.ExpectPing().WillReturnError(errors.New("database offline"))
			failedExpect.ExpectClose()
			pools=impl.NewPostgresPools(func(string)(*gorm.DB,error){return failed,nil})
			_,err=pools.Business(ctx)
			Expect(err).To(MatchError(ContainSubstring("ping postgres gorm db")))
			Expect(failedExpect.ExpectationsWereMet()).To(Succeed())
		})
		It("borrows the business pool, or opens and closes a separately owned Audit pool",func(){
			t:=GinkgoT()
			t.Setenv("DATABASE_URL","business")
			t.Setenv("AUDIT_DATABASE_URL","")
			t.Setenv("PLATFORM_DB_MAX_OPEN","")
			t.Setenv("PLATFORM_DB_MAX_IDLE","")
			t.Setenv("AUDIT_DB_MAX_OPEN","")
			t.Setenv("AUDIT_DB_MAX_IDLE","")
			ctx:=context.Background()
			borrowed:=&gorm.DB{Config:&gorm.Config{}}
			pools:=impl.NewPostgresPools(func(string)(*gorm.DB,error){Fail("borrowed pool must not open");return nil,nil})
			got,closePool,err:=pools.Audit(ctx,borrowed)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(BeIdenticalTo(borrowed))
			closePool()
			for _,dsn:=range []string{"","audit"}{
				t.Setenv("AUDIT_DATABASE_URL",dsn)
				db,expect:=mockPool()
				expect.ExpectPing()
				expect.ExpectClose()
				pools=impl.NewPostgresPools(func(got string)(*gorm.DB,error){
					want:="business"
					if dsn!=""{want=dsn}
					Expect(got).To(Equal(want))
					return db,nil
				})
				got,closePool,err:=pools.Audit(ctx,nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(got).To(BeIdenticalTo(db))
				raw,err:=got.DB()
				Expect(err).NotTo(HaveOccurred())
				if dsn==""{Expect(raw.Stats().MaxOpenConnections).To(Equal(20))}else{Expect(raw.Stats().MaxOpenConnections).To(Equal(5))}
				closePool()
				Expect(expect.ExpectationsWereMet()).To(Succeed())
			}
			pools=impl.NewPostgresPools(func(string)(*gorm.DB,error){return nil,errors.New("database unavailable")})
			got,closePool,err=pools.Audit(ctx,nil)
			Expect(got).To(BeNil())
			Expect(closePool).To(BeNil())
			Expect(err).To(MatchError("could not open Audit PostgreSQL database"))
		})
	})
}
