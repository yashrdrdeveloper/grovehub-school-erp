package impl

import (
 bookrepo "grovehub/services/accounting/internal/app/repository/accountingbook"
 bookimpl "grovehub/services/accounting/internal/app/repository/accountingbook/impl"
 entityrepo "grovehub/services/accounting/internal/app/repository/accountingentity"
 entityimpl "grovehub/services/accounting/internal/app/repository/accountingentity/impl"
)

type repositoryGraph struct {
 books bookrepo.Repository
 entities entityrepo.Repository
}

func (b *builder) newRepository() repositoryGraph {
 return repositoryGraph{books: bookimpl.New(b.db), entities: entityimpl.New(b.db)}
}
