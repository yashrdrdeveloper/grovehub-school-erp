package dto_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"encoding/json"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/dto"
)

var _ = Describe("UnmarshalJSON", func() {
	It("distinguishes omitted, null, zero and invalid values", func() {
		var value struct { Field dto.Patch[int] }
		Expect(json.Unmarshal([]byte("{}"), &value)).To(Succeed())
		Expect(value.Field.Set).To(BeFalse())
		Expect(json.Unmarshal([]byte(`{"Field":null}`), &value)).To(Succeed())
		Expect(value.Field.Set).To(BeTrue())
		Expect(value.Field.Value).To(BeNil())
		Expect(json.Unmarshal([]byte(`{"Field":0}`), &value)).To(Succeed())
		Expect(value.Field.Value).To(HaveValue(Equal(0)))
		Expect(value.Field.UnmarshalJSON([]byte(`"wrong"`))).To(HaveOccurred())
		Expect(value.Field.Value).To(BeNil())
	})
})
