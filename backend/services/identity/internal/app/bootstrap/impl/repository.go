package impl
import(
 tx "grovehub/libs/core/database/transaction"
 accountrepo "grovehub/services/identity/internal/app/repository/account"
 accountimpl "grovehub/services/identity/internal/app/repository/account/impl"
 credentialrepo "grovehub/services/identity/internal/app/repository/credential"
 credentialimpl "grovehub/services/identity/internal/app/repository/credential/impl"
 sessionrepo "grovehub/services/identity/internal/app/repository/session"
 sessionimpl "grovehub/services/identity/internal/app/repository/session/impl"
 tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
 tokenimpl "grovehub/services/identity/internal/app/repository/passwordtoken/impl"
 mailrepo "grovehub/services/identity/internal/app/repository/invitation"
 mailimpl "grovehub/services/identity/internal/app/repository/invitation/impl"
)
type repositories struct{accounts accountrepo.Repository;credentials credentialrepo.Repository;sessions sessionrepo.Repository;tokens tokenrepo.Repository;mail mailrepo.Repository;unit tx.UnitOfWork}
func(b *builder)newRepository()repositories{return repositories{
 accounts:accountimpl.New(b.db),credentials:credentialimpl.New(b.db),sessions:sessionimpl.New(b.db),tokens:tokenimpl.New(b.db),mail:mailimpl.New(b.db),unit:tx.NewUnitOfWork(b.db),
}}
