package impl

import (
 "context"
 "errors"
 "strings"

 "github.com/google/uuid"
 "grovehub/services/accounting/internal/app/dto"
 "grovehub/services/accounting/internal/app/entity"
 bookrepo "grovehub/services/accounting/internal/app/repository/accountingbook"
 entityrepo "grovehub/services/accounting/internal/app/repository/accountingentity"
 service "grovehub/services/accounting/internal/app/services/accountingbook"
 "grovehub/services/accounting/internal/app/validators"
)

type implementation struct {
 books bookrepo.Repository
 entities entityrepo.Repository
}

var _ service.Service = (*implementation)(nil)

func New(books bookrepo.Repository, entities entityrepo.Repository) service.Service {
 return &implementation{books: books, entities: entities}
}

func (s *implementation) Create(ctx context.Context, tenantID string, input dto.CreateAccountingBook) (entity.AccountingBook, error) {
 code, name := strings.ToUpper(strings.TrimSpace(input.Code)), strings.TrimSpace(input.Name)
 if err := validators.ValidateIdentity(tenantID, code, name); err != nil { return entity.AccountingBook{}, err }
 if err := validators.ValidateID(input.EntityID); err != nil { return entity.AccountingBook{}, err }
 currency := strings.ToUpper(strings.TrimSpace(input.BaseCurrency))
 if currency != "INR" { return entity.AccountingBook{}, errors.New("unsupported base currency") }
 if _, err := s.entities.FindByID(ctx, tenantID, input.EntityID); err != nil { return entity.AccountingBook{}, err }
 value := entity.AccountingBook{BookID: uuid.NewString(), TenantID: tenantID, EntityID: input.EntityID, Code: code, Name: name, BaseCurrency: currency, CurrencyScale: 2, AccountingBasis: "ACCRUAL"}
 if err := s.books.Create(ctx, value); err != nil { return entity.AccountingBook{}, err }
 return s.books.FindByID(ctx, tenantID, value.BookID)
}

func (s *implementation) FindByID(ctx context.Context, tenantID, id string) (entity.AccountingBook, error) {
 if err := validators.ValidateID(id); err != nil { return entity.AccountingBook{}, err }
 return s.books.FindByID(ctx, tenantID, id)
}
