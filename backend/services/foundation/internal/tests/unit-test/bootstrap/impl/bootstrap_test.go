package impl_test

import (
    "strings"
    "log/slog"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "gorm.io/gorm"

    bootstrap "grovehub/services/foundation/internal/app/bootstrap/impl"
)

// @Testing suite=FoundationBootstrap
func Specs(){
    Describe("Foundation bootstrap preconditions",func(){
        It("refuses to construct without a database",func(){
            builder:=bootstrap.New(slog.New(slog.NewTextHandler(&strings.Builder{},nil)),nil,nil)
            _,err:=builder.Build()
            Expect(err).To(MatchError("foundation bootstrap requires a database"))
        })
        It("refuses to construct without a durable event store",func(){
            builder:=bootstrap.New(nil,&gorm.DB{},nil)
            _,err:=builder.Build()
            Expect(err).To(MatchError("foundation bootstrap requires an event store"))
        })
    })
}
