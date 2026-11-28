package fixture

import (
	"bytes"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	"grovehub/services/identity/internal/tests/mocks/operatordatabase"
)

var Failure = errors.New("operator fixture failure")

type Writer struct{}

func (Writer) Write([]byte) (int, error) {
	return 0, Failure
}

// Operator owns a fresh mock database boundary and checks cleanup after each test.
type Operator struct {
	Mock        *database.MockDatabase
	BusinessDB  *gorm.DB
	AuditDB     *gorm.DB
	Output      bytes.Buffer
	Diagnostics bytes.Buffer
	Closed      []string
}

func New(t testing.TB) *Operator {
	return &Operator{
		Mock:       database.NewMockDatabase(gomock.NewController(t)),
		BusinessDB: &gorm.DB{Config: &gorm.Config{}},
		AuditDB:    &gorm.DB{Config: &gorm.Config{}},
	}
}

func (f *Operator) ExpectBusiness(err error) {
	if err != nil {
		f.Mock.EXPECT().Business(gomock.Any()).Return(nil, nil, err)
		return
	}

	f.Mock.EXPECT().Business(gomock.Any()).Return(f.BusinessDB, func() {
		f.Closed = append(f.Closed, "business")
	}, nil)
}

func (f *Operator) ExpectAudit(shared *gorm.DB, err error) {
	if err != nil {
		f.Mock.EXPECT().Audit(gomock.Any(), shared).Return(nil, nil, err)
		return
	}

	f.Mock.EXPECT().Audit(gomock.Any(), shared).Return(f.AuditDB, func() {
		f.Closed = append(f.Closed, "audit")
	}, nil)
}

func (f *Operator) ExpectName(db *gorm.DB, name string, err error) {
	f.Mock.EXPECT().CurrentName(gomock.Any(), db).Return(name, err)
}
