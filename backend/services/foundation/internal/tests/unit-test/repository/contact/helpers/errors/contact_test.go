package errors_test

import (
    "errors"
    "github.com/jackc/pgx/v5/pgconn"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    mapper "grovehub/services/foundation/internal/app/repository/contact/helpers/errors"
    repoerrors "grovehub/services/foundation/internal/app/repository/constants/errors"
    keys "grovehub/services/foundation/internal/app/repository/constants/keys"
)

// @Testing suite=ContactRepositoryErrorMapping
func ErrorSpecs() {
    Describe("MapContactWriteError",func(){
        It("preserves nil errors",func(){Expect(mapper.MapContactWriteError(nil)).To(Succeed())})
        It("maps contact uniqueness violations",func(){
            err:=&pgconn.PgError{Code:"23505",ConstraintName:keys.ContactTypeNormalizedValueUniqueConstraint}
            Expect(mapper.MapContactWriteError(err)).To(MatchError(repoerrors.ErrContactAlreadyExists))
        })
        It("preserves all unrelated constraint and ordinary errors",func(){
            for _,err:=range []error{
                &pgconn.PgError{Code:"23503",ConstraintName:keys.ContactTypeNormalizedValueUniqueConstraint},
                &pgconn.PgError{Code:"23505",ConstraintName:"another_constraint"},
                errors.New("failed query"),
            } {
                Expect(mapper.MapContactWriteError(err)).To(BeIdenticalTo(err))
            }
        })
    })
}
