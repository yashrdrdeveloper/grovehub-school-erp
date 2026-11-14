package helpers_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    coreerrors "grovehub/libs/core/errors"
    validators "grovehub/services/foundation/internal/app/validators/person/helpers"
    ids "grovehub/services/foundation/internal/tests/constants/person"
)

// @Testing suite=FoundationPersonValidation
func GIDSpecs() {
    Describe("ValidateGID",func(){
        It("accepts a valid UUIDv7",func(){Expect(validators.ValidateGID(ids.GID)).To(Succeed())})
        It("rejects an invalid identifier with a domain error",func(){
            err:=validators.ValidateGID("not-a-uuid")
            Expect(err).To(HaveOccurred())
            var typed *coreerrors.Error
            Expect(err).To(BeAssignableToTypeOf(typed))
        })
    })
}
