package constants_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/database/internal/src/postgres/operations/projectpool/constants"
	"time"
)

var _ = Describe("Pool limits", func() {
	It("preserves the configured connection budgets and lifetimes", func() {
		Expect(constants.MaxOpenConns).To(Equal(10))
		Expect(constants.MaxIdleConns).To(Equal(2))
		Expect(constants.ConnectionTimeout).To(Equal(5 * time.Second))
		Expect(constants.ConnMaxIdleTime).To(Equal(5 * time.Minute))
		Expect(constants.ConnMaxLifetime).To(Equal(30 * time.Minute))
	})
})
