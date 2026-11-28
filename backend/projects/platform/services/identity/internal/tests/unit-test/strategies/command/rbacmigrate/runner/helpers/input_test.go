package helpers_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/services/identity/internal/app/strategies/command/rbacmigrate/runner/helpers"
)

func TestMigrationInput(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Migration Input Suite")
}

var _ = Describe("ReadReviewed", func() {
	It("requires the exact reviewed bytes", func() {
		data := []byte(`{"tenantId":"school"}`)
		path := filepath.Join(GinkgoT().TempDir(), "plan.json")
		Expect(os.WriteFile(path, data, 0600)).To(Succeed())
		result, err := helpers.ReadReviewed(path, helpers.Digest(data))
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(data))
		_, err = helpers.ReadReviewed(path, strings.Repeat("0", 64))
		Expect(err).To(MatchError(helpers.ErrInvalidInput))
	})
})

var _ = Describe("Decode", func() {
	It("rejects unknown fields and trailing documents", func() {
		var result struct { TenantID string `json:"tenantId"` }
		Expect(helpers.Decode([]byte(`{"tenantId":"school"}`), &result)).To(Succeed())
		Expect(result.TenantID).To(Equal("school"))
		Expect(helpers.Decode([]byte(`{"unknown":1}`), &result)).To(MatchError(helpers.ErrInvalidInput))
		Expect(helpers.Decode([]byte(`{} {}`), &result)).To(MatchError(helpers.ErrInvalidInput))
	})
})
