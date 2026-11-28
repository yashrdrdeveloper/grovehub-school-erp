package validation_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/search/internal/src/contracts"
	"grovehub/libs/core/search/internal/src/operations/cursor"
	"grovehub/libs/core/search/internal/src/operations/validation"
)

var _ = Describe("Multiple sorts", func() {
	It("preserves priority and binds all cursor sort values", func() {
		allowed := map[string]contracts.CursorValueValidator{"status": validation.ValidateNonEmptyCursor, "createdAt": validation.ValidateTimeCursor}
		fallback := contracts.QuerySort{Field: "createdAt", Direction: contracts.DirectionDescending}
		request := []contracts.Sort{{Field: "status", Direction: contracts.DirectionAscending}, {Field: "createdAt", Direction: contracts.DirectionDescending}}
		rules, err := validation.NormalizeSorts(request, fallback, allowed, 4, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(rules[0].Field).To(Equal("status"))
		_, err = validation.NormalizeSorts([]contracts.Sort{request[0], request[0]}, fallback, allowed, 4, nil)
		Expect(err).To(HaveOccurred())
		token, err := cursor.EncodeCursor(contracts.CursorPayload{Field: "status", Direction: contracts.DirectionAscending, Value: "ACTIVE", CursorId: "id", SecondarySorts: request[1:], SecondaryValues: []string{"2026-11-18T00:00:00Z"}})
		Expect(err).NotTo(HaveOccurred())
		result, err := validation.NormalizeMultiCursor(token, rules, allowed, validation.ValidateNonEmptyCursor, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.SecondaryValues).To(HaveLen(1))
		rules[1].Direction = contracts.DirectionAscending
		_, err = validation.NormalizeMultiCursor(token, rules, allowed, validation.ValidateNonEmptyCursor, nil)
		Expect(err).To(HaveOccurred())
	})
})
