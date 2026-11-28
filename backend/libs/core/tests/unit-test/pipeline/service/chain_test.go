package service_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corepipeline "grovehub/libs/core/pipeline"
)

var _ = Describe("Chain", func() {
	It("stops at the first return", func() {
		state := 0
		value,created,err := corepipeline.RunService(context.Background(),&state,
			func(_ context.Context,s *int)corepipeline.Outcome[string]{*s++;return corepipeline.Next[string]()},
			func(_ context.Context,s *int)corepipeline.Outcome[string]{*s++;return corepipeline.Return("existing",false)},
			func(context.Context,*int)corepipeline.Outcome[string]{Fail("unexpected handler");return corepipeline.Next[string]()},
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("existing"))
		Expect(created).To(BeFalse())
		Expect(state).To(Equal(2))
	})

	It("rejects unresolved, empty and nil chains", func() {
		for _, handlers := range [][]corepipeline.Handler[int,string]{nil,{nil},{corepipeline.Continue[int,string]}} {
			_,_,err := corepipeline.RunService(context.Background(),new(int),handlers...)
			Expect(err).To(MatchError(corepipeline.ErrUnresolvedChain))
		}
	})

	It("propagates failure and requires recovery to return or fail", func() {
		sentinel:=errors.New("failure")
		_,_,err:=corepipeline.RunService(context.Background(),new(int),func(context.Context,*int)corepipeline.Outcome[string]{
			return corepipeline.RequireHandled(corepipeline.Next[string](),sentinel)
		})
		Expect(err).To(MatchError(sentinel))
		value,created,err:=corepipeline.RunService(context.Background(),new(int),func(context.Context,*int)corepipeline.Outcome[string]{
			return corepipeline.RequireHandled(corepipeline.Return("done",true),sentinel)
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("done"))
		Expect(created).To(BeTrue())
	})
})
