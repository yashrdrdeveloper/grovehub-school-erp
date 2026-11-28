package service_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	service "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline/internal/src/operations/service"
)

var _ = Describe("Decide",func(){
	It("chooses only the first matching rule and supports a fallback",func(){
		for _,state:=range []int{1,2}{
			value,_,err:=service.Run(context.Background(),&state,func(ctx context.Context,s *int)service.Outcome[string]{
				return service.Decide(ctx,s,service.Rule[int,string]{},
					service.When(func(s *int)bool{return *s==1},func(context.Context,*int)service.Outcome[string]{return service.Return("first",true)}),
					service.Otherwise(func(context.Context,*int)service.Outcome[string]{return service.Return("fallback",false)}),
				)
			})
			Expect(err).NotTo(HaveOccurred())
			if state==1 {Expect(value).To(Equal("first"))}else{Expect(value).To(Equal("fallback"))}
		}
		_,_,err:=service.Run(context.Background(),new(int),func(ctx context.Context,s *int)service.Outcome[string]{return service.Decide[int,string](ctx,s)})
		Expect(err).To(MatchError(service.ErrNoMatchingRule))
	})
})

var _ = Describe("Lookup",func(){
	It("distinguishes success, missing and unexpected errors",func(){
		missing:=errors.New("missing")
		unexpected:=errors.New("storage error")
		for _,item:=range []struct{err error;want string}{{nil,"found"},{missing,"not found"},{unexpected,"failed"}}{
			value,_,err:=service.Run(context.Background(),new(int),func(context.Context,*int)service.Outcome[string]{
				return service.Lookup(1,item.err,missing,func(int)service.Outcome[string]{return service.Return("found",false)},func()service.Outcome[string]{return service.Return("not found",false)},func(error)service.Outcome[string]{return service.Return("failed",false)})
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal(item.want))
		}
	})
})
