package unitofwork_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/operations/unitofwork"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/operations/unitofwork/impl"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/operations/run"
)

var _ = Describe("UnitOfWork", func() {
	It("exposes the Run contract independently of its concrete implementation", func() {
		var unit unitofwork.UnitOfWork = impl.NewUnitOfWork(nil)

		Expect(unit.Run(context.Background(), nil)).To(MatchError(run.ErrMissingDatabase))
	})
})
