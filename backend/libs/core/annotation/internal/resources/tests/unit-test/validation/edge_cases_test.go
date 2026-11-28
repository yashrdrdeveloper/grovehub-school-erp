package validation_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/validation"
	"grovehub/libs/core/annotation/internal/src/dto"
)

type EmbeddedName struct { Name string `normalize:"trim" validate:"notblank"` }

var _ = Describe("DTO validation edge cases",func(){
	It("preserves typed errors and stops at every failed processing stage",func(){
		Expect(validation.Error{Field:"Name",Rule:"required"}.Error()).To(Equal("Name violates required"))
		type invalidNormalize struct { Name string `normalize:"unknown"` }
		_,err:=validation.Normalize(invalidNormalize{Name:"x"})
		Expect(err).To(MatchError(`Name: unsupported normalizer "unknown"`))
		mapped:=errors.New("mapped normalization")
		_,err=validation.Process(invalidNormalize{},validation.Rules[invalidNormalize]{MapDTOError:func(error)error{return mapped}})
		Expect(err).To(MatchError(mapped))
		_,err=validation.Process(item{Name:"x"},validation.Rules[item]{Normalize:func(item)(item,error){return item{},mapped}})
		Expect(err).To(MatchError(mapped))
		_,err=validation.Process(item{Name:"x"},validation.Rules[item]{Validate:func(item)error{return mapped}})
		Expect(err).To(MatchError(mapped))
		_,err=validation.NormalizeFields(123)
		Expect(err).To(MatchError("normalization target must be a struct"))
		Expect(validation.Validate(123)).To(MatchError("validation target must be a struct"))
	})
	It("copies optional, embedded and nested values without touching unrelated fields",func(){
		type input struct {
			EmbeddedName
			ID string `normalize:"trim,uuidv7" validate:"uuidv7"`
			Ptr *string `normalize:"trim" validate:"notempty"`
			Nullable *string `normalize:"trim"`
			hidden string
			Untagged string
			Patch dto.Patch[string] `normalize:"trim" validate:"required"`
			Items []item `validate:"valid"`
		}
		text:=" a "
		value:=input{EmbeddedName:EmbeddedName{Name:" Ada "},ID:" 01890ABC-1234-7ABC-8DEF-0123456789AB ",Ptr:&text,hidden:" hidden ",Untagged:" untouched ",Items:[]item{{Name:" b "}}}
		out,err:=validation.Normalize(value)
		Expect(err).NotTo(HaveOccurred())
		Expect(out.Name).To(Equal("Ada"))
		Expect(out.ID).To(Equal("01890abc-1234-7abc-8def-0123456789ab"))
		Expect(out.Ptr).To(HaveValue(Equal("a")))
		Expect(text).To(Equal(" a "))
		Expect(out.hidden).To(Equal(" hidden "))
		Expect(out.Untagged).To(Equal(" untouched "))
		Expect(out.Nullable).To(BeNil())
		Expect(out.Items[0].Name).To(Equal("b"))
		Expect(value.Items[0].Name).To(Equal(" b "))
		value.ID="bad"
		_,err=validation.NormalizeFields(value)
		Expect(err).To(MatchError("ID violates uuidv7"))
		Expect(validation.Validate(value)).To(MatchError("ID violates uuidv7"))
	})
	It("rejects unsupported normalization targets and nested invalid normalization",func(){
		type invalidItem struct { Name string `normalize:"invalid"` }
		_,err:=validation.NormalizeFields(struct{N int `normalize:"trim"`}{})
		Expect(err).To(MatchError("N: unsupported normalize target"))
		_,err=validation.NormalizeFields(struct{N struct{} `normalize:"trim"`}{})
		Expect(err).To(MatchError("N: unsupported normalize target"))
		_,err=validation.NormalizeFields(struct{N []string `normalize:"trim"`}{})
		Expect(err).To(MatchError("N: unsupported normalize target"))
		_,err=validation.NormalizeFields(struct{Items []invalidItem `validate:"valid"`}{Items:[]invalidItem{{Name:"x"}}})
		Expect(err).To(MatchError(ContainSubstring("Items[0].Name: unsupported normalizer")))
	})
	It("validates string, pointer, collection and numeric rules",func(){
		type valid struct {
			Empty string `validate:"notempty"`
			Required string `validate:"required"`
			Nullable *string `validate:"nonnull"`
			Optional *string `validate:"notempty"`
			Unsigned uint16 `validate:"min=1,max=4"`
			Decimal float64 `validate:"min=0.5,max=2"`
			Names []string `validate:"notempty"`
			Mapping map[string]string `validate:"notempty"`
			Patch dto.Patch[string] `validate:"notblank"`
		}
		text:="x"
		v:=valid{Empty:" ",Required:"x",Nullable:&text,Unsigned:2,Decimal:1,Names:[]string{"x"},Mapping:map[string]string{"x":"y"},Patch:dto.Patch[string]{Set:true,Value:&text}}
		Expect(validation.Validate(v)).To(Succeed())
		for _,value:=range []struct{mutate func(*valid);message string}{
			{func(v *valid){v.Empty=""},"Empty violates notempty"},
			{func(v *valid){v.Required=" "},"Required violates required"},
			{func(v *valid){v.Nullable=nil},"Nullable violates nonnull"},
			{func(v *valid){v.Unsigned=0},"Unsigned violates min"},
			{func(v *valid){v.Decimal=3},"Decimal violates max"},
			{func(v *valid){v.Names=nil},"Names violates notempty"},
			{func(v *valid){v.Mapping=nil},"Mapping violates notempty"},
		}{copy:=v;value.mutate(&copy);Expect(validation.Validate(copy)).To(MatchError(value.message))}
	})
	It("rejects invalid rule declarations rather than silently accepting them",func(){
		for _,value:=range []struct{input any;message string}{
			{struct{Value string `validate:"maxrunes=-1"`}{}, "invalid maxrunes rule"},
			{struct{Value string `validate:"maxrunes=no"`}{}, "invalid maxrunes rule"},
			{struct{Value string `validate:"mystery"`}{}, "unsupported validation rule"},
			{struct{Value int `validate:"mystery"`}{}, "unsupported validation rule"},
			{struct{Value int `validate:"min=bad"`}{}, "invalid numeric constraint"},
			{struct{Value bool `validate:"min=1"`}{}, "numeric constraint on unsupported type"},
			{struct{Value int `validate:"valid"`}{}, "validate:valid requires a struct or slice"},
		}{Expect(validation.Validate(value.input)).To(MatchError(ContainSubstring(value.message)))}
	})
})
