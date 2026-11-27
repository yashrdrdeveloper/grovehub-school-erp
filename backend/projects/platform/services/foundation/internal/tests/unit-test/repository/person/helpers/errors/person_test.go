package errors_test

import (
    "errors"
    "github.com/jackc/pgx/v5/pgconn"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    mapper "grovehub/services/foundation/internal/app/repository/person/helpers/errors"
    repoerrors "grovehub/services/foundation/internal/app/repository/constants/errors"
    keys "grovehub/services/foundation/internal/app/repository/constants/keys"
)

// @Testing suite=PersonRepositoryErrorMapping
func ErrorSpecs() {
    Describe("MapPersonWriteError",func(){
        It("preserves nil errors",func(){Expect(mapper.MapPersonWriteError(nil)).To(Succeed())})
        It("maps known uniqueness violations",func(){
            err:=&pgconn.PgError{Code:"23505",ConstraintName:keys.PersonIdempotencyKeyUniqueConstraint}
            Expect(mapper.MapPersonWriteError(err)).To(MatchError(repoerrors.ErrPersonIdempotencyKeyAlreadyExists))
        })
        It("does not rewrite unrelated database errors",func(){
            err:=&pgconn.PgError{Code:"23503",ConstraintName:keys.PersonIdempotencyKeyUniqueConstraint}
            Expect(mapper.MapPersonWriteError(err)).To(BeIdenticalTo(err))
        })
        It("does not rewrite an unrelated unique constraint",func(){
            err:=&pgconn.PgError{Code:"23505",ConstraintName:"different"}
            Expect(mapper.MapPersonWriteError(err)).To(BeIdenticalTo(err))
        })
        It("preserves ordinary errors",func(){
            err:=errors.New("query failure")
            Expect(mapper.MapPersonWriteError(err)).To(BeIdenticalTo(err))
        })
    })
}
