// persistence owns the annotation definitions for this backend responsibility.
package persistence

import "grovehub/libs/core/annotation/contract"

var UseCases = []contract.UseCase{
    {Name:"Repository",Syntax:"@Repository",Category:"Persistence, repository, transactions and search",Purpose:"Repository implementation metadata and constructor wiring",Mode:"REG",Status:contract.Proposed,Roadmap:"P4",Target:contract.InterfaceTarget},
    {Name:"CRUD",Syntax:"@CRUD",Category:"Persistence, repository, transactions and search",Purpose:"Opt-in GORM basic CRUD only; no guessed aggregate persistence",Mode:"GEN",Status:contract.Optional,Roadmap:"P4",Target:contract.InterfaceTarget},
    {Name:"Query",Syntax:"@Query",Category:"Persistence, repository, transactions and search",Purpose:"Explicit typed query declaration; parameterized only",Mode:"GEN/CHECK",Status:contract.Optional,Roadmap:"P4",Target:contract.MethodTarget},
    {Name:"FindBy",Syntax:"@FindBy",Category:"Persistence, repository, transactions and search",Purpose:"Allowlisted trivial finder generation; fail on ambiguity",Mode:"GEN",Status:contract.Optional,Roadmap:"P4",Target:contract.MethodTarget},
    {Name:"Transactional",Syntax:"@Transactional",Category:"Persistence, repository, transactions and search",Purpose:"Generated wrapper using core/database/transaction",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P4",Target:contract.MethodTarget},
    {Name:"ReadOnly",Syntax:"@ReadOnly",Category:"Persistence, repository, transactions and search",Purpose:"Read-only transaction intent; enforce at DB where possible",Mode:"WRAP",Status:contract.Optional,Roadmap:"P4",Target:contract.MethodTarget},
    {Name:"OptimisticLock",Syntax:"@OptimisticLock",Category:"Persistence, repository, transactions and search",Purpose:"Version predicate plus conflict mapping",Mode:"GEN",Status:contract.Optional,Roadmap:"P4",Target:contract.MethodTarget},
    {Name:"TenantScoped",Syntax:"@TenantScoped",Category:"Persistence, repository, transactions and search",Purpose:"Require trusted tenant context in repository predicates",Mode:"CHECK/WRAP",Status:contract.Proposed,Roadmap:"P7",Target:contract.InterfaceTarget},
    {Name:"SoftDelete",Syntax:"@SoftDelete",Category:"Persistence, repository, transactions and search",Purpose:"Optional GORM soft-delete integration",Mode:"GEN",Status:contract.Optional,Roadmap:"P4",Target:contract.InterfaceTarget},
    {Name:"Searchable",Syntax:"@Searchable",Category:"Persistence, repository, transactions and search",Purpose:"Declare allowed search fields, not raw SQL expressions",Mode:"REG/CHECK",Status:contract.Proposed,Roadmap:"P4",Target:contract.MethodTarget},
    {Name:"Filterable",Syntax:"@Filterable",Category:"Persistence, repository, transactions and search",Purpose:"Allowed filters/operators with canonicalization",Mode:"REG/CHECK",Status:contract.Optional,Roadmap:"P4",Target:contract.MethodTarget},
    {Name:"Sortable",Syntax:"@Sortable",Category:"Persistence, repository, transactions and search",Purpose:"Allowed stable sort keys and null ordering",Mode:"REG/CHECK",Status:contract.Optional,Roadmap:"P4",Target:contract.MethodTarget},
    {Name:"Cursor",Syntax:"@Cursor",Category:"Persistence, repository, transactions and search",Purpose:"Typed cursor serialization/validation policy",Mode:"GEN/CHECK",Status:contract.Optional,Roadmap:"P4",Target:contract.MethodTarget},
    {Name:"Migration",Syntax:"@Migration",Category:"Persistence, repository, transactions and search",Purpose:"Only mark migration files for linting; never auto-apply DDL",Mode:"CHECK",Status:contract.Discouraged,Roadmap:"P4",Target:contract.MethodTarget},
}
