package contracts_test

import (
	"encoding/json"
	coreerrors "grovehub/libs/core/errors"
	"grovehub/libs/core/http/internal/src/contracts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewSuccess", func() {
	It("constructs the existing JSON envelope", func() {
		raw, err := json.Marshal(contracts.NewSuccess("data"))
		Expect(err).NotTo(HaveOccurred())
		Expect(string(raw)).To(MatchJSON(`{"success":true,"data":"data"}`))
	})
})

var _ = Describe("NewError", func() {
	It("constructs the existing error envelope with typed metadata", func() {
		raw, err := json.Marshal(contracts.NewError(coreerrors.Code("INVALID"), "Invalid", 400, map[string]string{"field":"name"}))
		Expect(err).NotTo(HaveOccurred())
		Expect(string(raw)).To(MatchJSON(`{"success":false,"errorCode":"INVALID","errorMessage":"Invalid","httpErrorCode":400,"errorMetaData":{"field":"name"}}`))
	})
})
