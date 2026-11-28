package service_test

import (
	"context"
	"errors"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline"
)

var _ = Describe("Pipeline facade", func() {
	It("exposes result creation and direct chain execution", func() {
		Expect(corepipeline.Created("new").StatusCode).To(Equal(http.StatusCreated))
		Expect(corepipeline.CreatedOrOK("new",true)).To(Equal(corepipeline.Created("new")))
		Expect(corepipeline.CreatedOrOK("old",false)).To(Equal(corepipeline.OK("old")))
		state := 0
		outcome := corepipeline.Chain(context.Background(), &state, func(context.Context,*int)corepipeline.Outcome[string] {
			return corepipeline.Return("done",true)
		})
		value,created,err := corepipeline.RunService(context.Background(),&state,func(context.Context,*int)corepipeline.Outcome[string] { return outcome })
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("done"))
		Expect(created).To(BeTrue())
	})

	It("exposes typed error-first result resolution",func(){
		sentinel:=errors.New("failed")
		for _,input:=range []error{nil,sentinel}{
			value,_,err:=corepipeline.RunService(context.Background(),new(int),func(context.Context,*int)corepipeline.Outcome[string]{
				return corepipeline.ResolveResult(1,input,func(int)corepipeline.Outcome[string]{return corepipeline.Return("ok",false)},func(err error)corepipeline.Outcome[string]{return corepipeline.Fail[string](err)})
			})
			if input==nil {Expect(err).NotTo(HaveOccurred());Expect(value).To(Equal("ok"))}else{Expect(err).To(MatchError(sentinel))}
		}
	})
})
