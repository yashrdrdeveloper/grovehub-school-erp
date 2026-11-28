package authorization_test

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	mockdatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/tests/mocks/database"
	schoolgateway "grovehub/launcher/services/gateway/session"
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/entity"
	authorizationrepo "grovehub/services/identity/internal/app/repository/authorization"
	"grovehub/services/identity/internal/app/repository/authorization/helpers"
	authorizationimpl "grovehub/services/identity/internal/app/repository/authorization/impl"
	"grovehub/services/identity/internal/app/services/authorization/operations/evaluate"
)

const accountID = "11111111-1111-4111-8111-111111111111"
const tenantID = "school-1"

func snapshotRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"account_id", "tenant_id", "account_status", "role_id", "role_tenant_id", "role_status", "permission_code", "permission_active", "feature_code", "feature_active", "feature_enabled"})
}

var _ = Describe("LoadSnapshot", func() {
	var repository authorizationrepo.Repository
	var mock sqlmock.Sqlmock
	var query string
	var ctx context.Context

	BeforeEach(func() {
		db, driver, closeDB, err := mockdatabase.NewPostgres()
		Expect(err).NotTo(HaveOccurred())
		mock = driver
		DeferCleanup(func() {
			mock.ExpectClose()
			Expect(closeDB()).To(Succeed())
			Expect(mock.ExpectationsWereMet()).To(Succeed())
		})
		repository = authorizationimpl.New(db)
		query = regexp.QuoteMeta(strings.Replace(strings.Replace(helpers.SnapshotQuery, "?", "$1", 1), "?", "$2", 1))
		ctx = schoolgateway.WithPrincipal(context.Background(), schoolgateway.Principal{AccountID: accountID, TenantID: tenantID, SessionID: "session-1"})
	})

	It("loads grants and feature states from one tenant-scoped statement", func() {
		rows := snapshotRows().
			AddRow(accountID, tenantID, string(entity.Active), "role-1", tenantID, "ACTIVE", string(permission.AccountSearch), true, string(feature.AccountManagement), true, true).
			AddRow(accountID, tenantID, string(entity.Active), "role-1", tenantID, "ACTIVE", string(permission.AccountSelfView), true, nil, nil, nil).
			AddRow(accountID, tenantID, string(entity.Active), "role-2", tenantID, "INACTIVE", string(permission.AccountSearch), true, string(feature.AccountManagement), true, true)
		mock.ExpectQuery(query).WithArgs(tenantID, accountID).WillReturnRows(rows)

		snapshot, err := repository.LoadSnapshot(ctx, tenantID, accountID)
		Expect(err).NotTo(HaveOccurred())
		Expect(snapshot.Roles).To(HaveLen(2))
		Expect(snapshot.Roles[0].Permissions).To(ConsistOf(permission.AccountSearch, permission.AccountSelfView))
		Expect(snapshot.Roles[1].Active).To(BeFalse())
		Expect(evaluate.Allowed(ctx, permission.AccountSearch, snapshot)).To(BeTrue())
		Expect(evaluate.Allowed(ctx, permission.AccountSelfView, snapshot)).To(BeTrue())
		Expect(evaluate.Allowed(ctx, permission.AccountDeactivate, snapshot)).To(BeFalse())
	})

	It("denies a permission if the tenant entitlement is absent", func() {
		rows := snapshotRows().AddRow(accountID, tenantID, string(entity.Active), "role-1", tenantID, "ACTIVE", string(permission.AccountSearch), true, string(feature.AccountManagement), true, nil)
		mock.ExpectQuery(query).WithArgs(tenantID, accountID).WillReturnRows(rows)

		snapshot, err := repository.LoadSnapshot(ctx, tenantID, accountID)
		Expect(err).NotTo(HaveOccurred())
		Expect(evaluate.Allowed(ctx, permission.AccountSearch, snapshot)).To(BeFalse())
	})

	It("loads an account without grants without synthesizing a role", func() {
		rows := snapshotRows().AddRow(accountID, tenantID, string(entity.Active), nil, nil, nil, nil, nil, nil, nil, nil)
		mock.ExpectQuery(query).WithArgs(tenantID, accountID).WillReturnRows(rows)

		snapshot, err := repository.LoadSnapshot(ctx, tenantID, accountID)
		Expect(err).NotTo(HaveOccurred())
		Expect(snapshot.Roles).To(BeEmpty())
		Expect(evaluate.Allowed(ctx, permission.AccountSearch, snapshot)).To(BeFalse())
	})

	It("rejects missing scope before querying", func() {
		_, err := repository.LoadSnapshot(ctx, "", accountID)
		Expect(err).To(MatchError(authorizationrepo.ErrScopeRequired))
		_, err = repository.LoadSnapshot(ctx, tenantID, "")
		Expect(err).To(MatchError(authorizationrepo.ErrScopeRequired))
	})

	It("returns not found for accounts outside the requested scope", func() {
		mock.ExpectQuery(query).WithArgs(tenantID, accountID).WillReturnRows(snapshotRows())

		snapshot, err := repository.LoadSnapshot(ctx, tenantID, accountID)
		Expect(err).To(MatchError(authorizationrepo.ErrNotFound))
		Expect(snapshot).To(Equal(authorizationrepo.Snapshot{}))
	})

	It("rejects inconsistent returned tenant or role scope", func() {
		rows := snapshotRows().AddRow(accountID, "school-2", string(entity.Active), nil, nil, nil, nil, nil, nil, nil, nil)
		mock.ExpectQuery(query).WithArgs(tenantID, accountID).WillReturnRows(rows)
		_, err := repository.LoadSnapshot(ctx, tenantID, accountID)
		Expect(err).To(MatchError(authorizationrepo.ErrInvalidSnapshot))

		rows = snapshotRows().AddRow(accountID, tenantID, string(entity.Active), "role-1", "school-2", "ACTIVE", nil, nil, nil, nil, nil)
		mock.ExpectQuery(query).WithArgs(tenantID, accountID).WillReturnRows(rows)
		_, err = repository.LoadSnapshot(ctx, tenantID, accountID)
		Expect(err).To(MatchError(authorizationrepo.ErrInvalidSnapshot))
	})

	It("propagates database errors without grants", func() {
		dbErr := errors.New("database unavailable")
		mock.ExpectQuery(query).WithArgs(tenantID, accountID).WillReturnError(dbErr)

		snapshot, err := repository.LoadSnapshot(ctx, tenantID, accountID)
		Expect(err).To(MatchError(dbErr))
		Expect(snapshot).To(Equal(authorizationrepo.Snapshot{}))
	})
})
