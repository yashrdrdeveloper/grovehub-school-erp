package response_test

import (
	"context"
	"net/http/httptest"
	coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/internal/src/operations/response"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("WriteSuccess", func() {
	It("writes the existing JSON contract and status", func() {
		w := httptest.NewRecorder()
		Expect(response.WriteSuccess(w,201,"ok")).To(Succeed())
		Expect(w.Code).To(Equal(201))
		Expect(w.Header().Get("Content-Type")).To(Equal("application/json; charset=utf-8"))
		Expect(w.Body.String()).To(MatchJSON(`{"success":true,"data":"ok"}`))
	})
})

var _ = Describe("WriteErrorWithMetadataContext", func() {
	It("preserves the error contract without requiring an active span", func() {
		w := httptest.NewRecorder()
		definition := coreerrors.NewDefinition("INVALID","Invalid",400)
		Expect(response.WriteErrorWithMetadataContext(context.Background(),w,400,definition,map[string]string{"field":"name"})).To(Succeed())
		Expect(w.Code).To(Equal(400))
		Expect(w.Body.String()).To(MatchJSON(`{"success":false,"errorCode":"INVALID","errorMessage":"Invalid","httpErrorCode":400,"errorMetaData":{"field":"name"}}`))
	})
})
