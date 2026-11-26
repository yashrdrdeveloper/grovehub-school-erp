package impl

import (
 "context"
 "strings"
 "time"

 "github.com/google/uuid"
 "grovehub/services/accounting/internal/app/dto"
 "grovehub/services/accounting/internal/app/entity"
 repo "grovehub/services/accounting/internal/app/repository/accountingentity"
 service "grovehub/services/accounting/internal/app/services/accountingentity"
 "grovehub/services/accounting/internal/app/validators"
)

type implementation struct { repository repo.Repository }

var _ service.Service = (*implementation)(nil)

func New(repository repo.Repository) service.Service {
 return &implementation{repository: repository}
}

func (s *implementation) Create(ctx context.Context, tenantID string, input dto.CreateAccountingEntity) (entity.AccountingEntity, error) {
 code, name := strings.ToUpper(strings.TrimSpace(input.Code)), strings.TrimSpace(input.Name)
 if err := validators.ValidateIdentity(tenantID, code, name); err != nil { return entity.AccountingEntity{}, err }
 value := entity.AccountingEntity{EntityID: uuid.NewString(), TenantID: tenantID, Code: code, Name: name, CreatedAt: time.Now().UTC()}
 if err := s.repository.Create(ctx, value); err != nil { return entity.AccountingEntity{}, err }
 return value, nil
}

func (s *implementation) FindByID(ctx context.Context, tenantID, id string) (entity.AccountingEntity, error) {
 if err := validators.ValidateID(id); err != nil { return entity.AccountingEntity{}, err }
 return s.repository.FindByID(ctx, tenantID, id)
}
