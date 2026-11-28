package validation_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/search/internal/src/contracts"
	"grovehub/libs/core/search/internal/src/operations/cursor"
	"grovehub/libs/core/search/internal/src/operations/validation"
)

var _ = Describe("Validation failure paths", func() {
	It("preserves resource-owned errors and rejects unsupported filter shapes", func() {
		invalid := errors.New("resource invalid")
		_, err := validation.NormalizeSearchTerms("long", 1, 1, invalid)
		Expect(err).To(MatchError(invalid))
		fields := map[string]contracts.FieldSpec{"name": {Operators: validation.EqualityOperators(), Normalizer: func(_ contracts.Operator, value string) (string, error) { return value, nil }}}
		for _, filter := range []contracts.Filter{{Field: "name", Operator: contracts.OperatorExists}, {Field: "name", Operator: contracts.OperatorEquals}, {Field: "name", Operator: contracts.OperatorIn, Values: []string{"a", "b"}}} {
			_, err := validation.NormalizeFilters([]contracts.Filter{filter}, fields, 1, 1, invalid)
			Expect(err).To(MatchError(invalid))
		}
		_, err = validation.NormalizeFilters(nil, fields, -1, 1, invalid)
		Expect(err).To(MatchError(invalid))
		fallback := contracts.QuerySort{Field: "name", Direction: contracts.DirectionAscending}
		_, err = validation.NormalizeSort([]contracts.Sort{{Field: "name"}, {Field: "name"}}, fallback, nil, invalid)
		Expect(err).To(MatchError(invalid))
	})
	It("handles default sorts, invalid limits and unknown fields", func() {
		fallback := contracts.QuerySort{Field: "name", Direction: contracts.DirectionAscending}
		allowed := map[string]contracts.CursorValueValidator{"name": validation.ValidateNonEmptyCursor}
		rules, err := validation.NormalizeSorts(nil, fallback, allowed, 2, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(rules).To(Equal([]contracts.QuerySort{fallback}))
		_, err = validation.NormalizeSorts(nil, fallback, nil, 2, nil)
		Expect(err).To(MatchError(contracts.ErrInvalidQuery))
		_, err = validation.NormalizeSorts([]contracts.Sort{{Field: "name"}}, fallback, allowed, 0, nil)
		Expect(err).To(MatchError(contracts.ErrInvalidQuery))
		_, err = validation.NormalizeSorts([]contracts.Sort{{Field: "unknown"}}, fallback, allowed, 2, nil)
		Expect(err).To(MatchError(contracts.ErrInvalidQuery))
	})
	It("validates optional cursors, malformed tokens and invalid decoded values", func() {
		rules := []contracts.QuerySort{{Field: "name", Direction: contracts.DirectionAscending}}
		allowed := map[string]contracts.CursorValueValidator{"name": validation.ValidateNonEmptyCursor}
		_, err := validation.NormalizeMultiCursor("", nil, allowed, validation.ValidateNonEmptyCursor, nil)
		Expect(err).To(MatchError(contracts.ErrInvalidQuery))
		_, err = validation.NormalizeMultiCursor("", rules, nil, validation.ValidateNonEmptyCursor, nil)
		Expect(err).To(MatchError(contracts.ErrInvalidQuery))
		result, err := validation.NormalizeMultiCursor(" ", rules, allowed, validation.ValidateNonEmptyCursor, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeNil())
		_, err = validation.NormalizeMultiCursor("!", rules, allowed, validation.ValidateNonEmptyCursor, nil)
		Expect(err).To(MatchError(contracts.ErrInvalidQuery))
		token, err := cursor.EncodeCursor(contracts.CursorPayload{Field: "name", Direction: contracts.DirectionAscending, Value: "", CursorId: "id"})
		Expect(err).NotTo(HaveOccurred())
		_, err = validation.NormalizeMultiCursor(token, rules, allowed, validation.ValidateNonEmptyCursor, nil)
		Expect(err).To(MatchError(contracts.ErrInvalidQuery))
		_, _, err = validation.NormalizeCursor(token, rules[0], validation.ValidateNonEmptyCursor, validation.ValidateNonEmptyCursor, nil)
		Expect(err).To(MatchError(contracts.ErrInvalidQuery))
	})
})
