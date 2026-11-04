package auth_test

import (
    "bytes"
    "net/http"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    httprequest "grovehub/libs/core/http/request"
    identityrequest "grovehub/services/identity/dto/request"
)

var _ = Describe("Register JSON contract", func() {
    It("rejects a caller-supplied gid", func() {
        req, err := http.NewRequest(http.MethodPost, "/auth/register",
            bytes.NewBufferString(`{"gid":"01900000-0000-7000-8000-000000000000","firstName":"Arun","middleName":"","lastName":"Kumar","email":"arun@example.com","password":"password123"}`))
        Expect(err).NotTo(HaveOccurred())
        var input identityrequest.Register
        Expect(httprequest.DecodeStrictJSON(req, &input)).To(MatchError(ContainSubstring("unknown field")))
    })

    It("accepts the register request without gid", func() {
        req, err := http.NewRequest(http.MethodPost, "/auth/register",
            bytes.NewBufferString(`{"firstName":"Arun","middleName":"","lastName":"Kumar","email":"arun@example.com","password":"password123"}`))
        Expect(err).NotTo(HaveOccurred())
        var input identityrequest.Register
        Expect(httprequest.DecodeStrictJSON(req, &input)).To(Succeed())
        Expect(input.Email).To(Equal("arun@example.com"))
    })
})
