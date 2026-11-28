package validation_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"errors"
	"grovehub/libs/core/annotation/internal/src/validation"
	 "grovehub/libs/core/annotation/internal/src/dto"
)

type item struct { Name string `normalize:"trim" validate:"notblank,maxrunes=3"` }
type request struct { Items []item `validate:"valid,notempty"`; Count int `validate:"min=1,max=3"`; Patch dto.Patch[string] `normalize:"trim" validate:"nonnull"` }

var _ = Describe("NormalizeFields", func() {
	It("copies slices and present PATCH pointers before normalizing", func() {
		text := " x "
		input := request{Items:[]item{{Name:" a "}}, Count:1, Patch:dto.Patch[string]{Set:true,Value:&text}}
		output,err:=validation.NormalizeFields(input)
		Expect(err).NotTo(HaveOccurred())
		Expect(output.Items[0].Name).To(Equal("a"))
		Expect(output.Patch.Value).To(HaveValue(Equal("x")))
		Expect(input.Items[0].Name).To(Equal(" a "))
		Expect(text).To(Equal(" x "))
	})
})

var _ = Describe("Validate", func() {
	It("reports precise nested field paths and rule failures", func() {
		Expect(validation.Validate(request{Items:[]item{{Name:""}},Count:1})).To(MatchError("Items[0].Name violates notblank"))
		Expect(validation.Validate(request{Items:[]item{{Name:"四五六七"}},Count:1})).To(MatchError("Items[0].Name violates maxrunes"))
		Expect(validation.Validate(request{Items:[]item{{Name:"a"}},Count:4})).To(MatchError("Count violates max"))
		Expect(validation.Validate(request{Items:[]item{{Name:"a"}},Count:1,Patch:dto.Patch[string]{Set:true}})).To(MatchError("Patch violates nonnull"))
	})
})

var _ = Describe("Process", func() {
	It("stops before custom validation when DTO validation fails and maps the error", func() {
		mapped:=errors.New("mapped")
		called:=false
		_,err:=validation.Process(item{},validation.Rules[item]{MapDTOError:func(error)error{return mapped},Validate:func(item)error{called=true;return nil}})
		Expect(err).To(MatchError(mapped))
		Expect(called).To(BeFalse())
	})
})
