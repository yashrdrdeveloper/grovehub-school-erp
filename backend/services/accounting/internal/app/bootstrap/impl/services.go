package impl

import (
 "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
 "grovehub/services/audit"
 auditclient "grovehub/services/audit/client/producer"
 bookservice "grovehub/services/accounting/internal/app/services/accountingbook"
 bookimpl "grovehub/services/accounting/internal/app/services/accountingbook/impl"
 entityservice "grovehub/services/accounting/internal/app/services/accountingentity"
 entityimpl "grovehub/services/accounting/internal/app/services/accountingentity/impl"
)

type serviceGraph struct {
 books bookservice.Service
 entities entityservice.Service
 unitOfWork transaction.UnitOfWork
 audit auditclient.Client
}

func (b *builder) newServices(repositories repositoryGraph) serviceGraph {
 return serviceGraph{
  books: bookimpl.New(repositories.books, repositories.entities),
  entities: entityimpl.New(repositories.entities),
  unitOfWork: transaction.NewUnitOfWork(b.db),
  audit: audit.NewProducer(b.events),
 }
}
