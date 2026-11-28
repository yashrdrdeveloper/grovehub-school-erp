package service_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	service "grovehub/libs/core/pipeline/internal/src/operations/service"
)

var _ = Describe("Chain", func() {
	It("stops at the first return", func() {
		state := 0
		value,created,err := service.Run(context.Background(),&state,
			func(_ context.Context,s *int)service.Outcome[string]{*s++;return service.Next[string]()},
			func(_ context.Context,s *int)service.Outcome[string]{*s++;return service.Return("existing",false)},
			func(context.Context,*int)service.Outcome[string]{Fail("unexpected handler");return service.Next[string]()},
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("existing"))
		Expect(created).To(BeFalse())
		Expect(state).To(Equal(2))
	})

	It("rejects unresolved, empty and nil chains", func() {
		for _, handlers := range [][]service.Handler[int,string]{nil,{nil},{service.Continue[int,string]}} {
			_,_,err := service.Run(context.Background(),new(int),handlers...)
			Expect(err).To(MatchError(service.ErrUnresolvedChain))
		}
	})

	It("propagates failure and requires recovery to return or fail", func() {
		sentinel:=errors.New("failure")
		_,_,err:=service.Run(context.Background(),new(int),func(context.Context,*int)service.Outcome[string]{
			return service.RequireHandled(service.Next[string](),sentinel)
		})
		Expect(err).To(MatchError(sentinel))
		value,created,err:=service.Run(context.Background(),new(int),func(context.Context,*int)service.Outcome[string]{
			return service.RequireHandled(service.Return("done",true),sentinel)
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("done"))
		Expect(created).To(BeTrue())
	})
})
