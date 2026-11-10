// Phase 1 persistence, repository, transactions and search annotation metadata. No runtime behavior.
package annotation

var persistenceUseCases = []UseCase{
    {Name:"Repository",Syntax:"@Repository",Category:"Persistence, repository, transactions and search",Purpose:"Repository implementation metadata and constructor wiring",Mode:"REG",Status:Proposed,Roadmap:"P4",Target:InterfaceTarget},
    {Name:"CRUD",Syntax:"@CRUD",Category:"Persistence, repository, transactions and search",Purpose:"Opt-in GORM basic CRUD only; no guessed aggregate persistence",Mode:"GEN",Status:Optional,Roadmap:"P4",Target:InterfaceTarget},
    {Name:"Query",Syntax:"@Query",Category:"Persistence, repository, transactions and search",Purpose:"Explicit typed query declaration; parameterized only",Mode:"GEN/CHECK",Status:Optional,Roadmap:"P4",Target:MethodTarget},
    {Name:"FindBy",Syntax:"@FindBy",Category:"Persistence, repository, transactions and search",Purpose:"Allowlisted trivial finder generation; fail on ambiguity",Mode:"GEN",Status:Optional,Roadmap:"P4",Target:MethodTarget},
    {Name:"Transactional",Syntax:"@Transactional",Category:"Persistence, repository, transactions and search",Purpose:"Generated wrapper using core/database/transaction",Mode:"WRAP",Status:Proposed,Roadmap:"P4",Target:MethodTarget},
    {Name:"ReadOnly",Syntax:"@ReadOnly",Category:"Persistence, repository, transactions and search",Purpose:"Read-only transaction intent; enforce at DB where possible",Mode:"WRAP",Status:Optional,Roadmap:"P4",Target:MethodTarget},
    {Name:"OptimisticLock",Syntax:"@OptimisticLock",Category:"Persistence, repository, transactions and search",Purpose:"Version predicate plus conflict mapping",Mode:"GEN",Status:Optional,Roadmap:"P4",Target:MethodTarget},
    {Name:"TenantScoped",Syntax:"@TenantScoped",Category:"Persistence, repository, transactions and search",Purpose:"Require trusted tenant context in repository predicates",Mode:"CHECK/WRAP",Status:Proposed,Roadmap:"P7",Target:InterfaceTarget},
    {Name:"SoftDelete",Syntax:"@SoftDelete",Category:"Persistence, repository, transactions and search",Purpose:"Optional GORM soft-delete integration",Mode:"GEN",Status:Optional,Roadmap:"P4",Target:InterfaceTarget},
    {Name:"Searchable",Syntax:"@Searchable",Category:"Persistence, repository, transactions and search",Purpose:"Declare allowed search fields, not raw SQL expressions",Mode:"REG/CHECK",Status:Proposed,Roadmap:"P4",Target:MethodTarget},
    {Name:"Filterable",Syntax:"@Filterable",Category:"Persistence, repository, transactions and search",Purpose:"Allowed filters/operators with canonicalization",Mode:"REG/CHECK",Status:Optional,Roadmap:"P4",Target:MethodTarget},
    {Name:"Sortable",Syntax:"@Sortable",Category:"Persistence, repository, transactions and search",Purpose:"Allowed stable sort keys and null ordering",Mode:"REG/CHECK",Status:Optional,Roadmap:"P4",Target:MethodTarget},
    {Name:"Cursor",Syntax:"@Cursor",Category:"Persistence, repository, transactions and search",Purpose:"Typed cursor serialization/validation policy",Mode:"GEN/CHECK",Status:Optional,Roadmap:"P4",Target:MethodTarget},
    {Name:"Migration",Syntax:"@Migration",Category:"Persistence, repository, transactions and search",Purpose:"Only mark migration files for linting; never auto-apply DDL",Mode:"CHECK",Status:Discouraged,Roadmap:"P4",Target:MethodTarget},
}
