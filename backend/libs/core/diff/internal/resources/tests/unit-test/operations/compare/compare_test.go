package compare_test

import (
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/diff/internal/src/contracts"
	"grovehub/libs/core/diff/internal/src/operations/compare"
)

type failingValue struct {
	fail bool
}

var marshalFailure = errors.New("test marshal failure")

func (value failingValue) MarshalJSON() ([]byte, error) {
	if value.fail {
		return nil, marshalFailure
	}

	return []byte("1"), nil
}

var _ = Describe("Compare", func() {
	It("wraps the original error when the old value cannot be marshaled", func() {
		changes, err := compare.Compare(failingValue{fail: true}, failingValue{})

		Expect(changes).To(BeNil())
		Expect(err).To(MatchError(marshalFailure))
		Expect(err.Error()).To(HavePrefix("marshal old value:"))
	})

	It("wraps the original error when the new value cannot be marshaled", func() {
		changes, err := compare.Compare(failingValue{}, failingValue{fail: true})

		Expect(changes).To(BeNil())
		Expect(err).To(MatchError(marshalFailure))
		Expect(err.Error()).To(HavePrefix("marshal new value:"))
	})

	It("preserves precise large integer values in changes", func() {
		changes, err := compare.Compare(uint64(9007199254740992), uint64(9007199254740993))

		Expect(err).NotTo(HaveOccurred())
		Expect(changes).To(Equal([]contracts.Change{{
			Field: "$", OldValue: json.Number("9007199254740992"), NewValue: json.Number("9007199254740993"),
		}}))
	})

	It("returns an allocated empty slice for equal values", func() {
		changes, err := compare.Compare("same", "same")

		Expect(err).NotTo(HaveOccurred())
		Expect(changes).NotTo(BeNil())
		Expect(changes).To(BeEmpty())
	})
})
