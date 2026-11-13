// persistence owns the annotation definitions for this backend responsibility.
package persistence

import "grovehub/libs/core/annotation/contract"

var UseCases = []contract.UseCase{
    {Name:"Repository",Syntax:"@Repository entity=Type",Category:"Persistence, repository, transactions and search",Purpose:"Validated repository ownership metadata; no constructor or transaction is generated",Mode:"CHECK",Status:contract.Pilot,Roadmap:"P4",Target:contract.InterfaceTarget},
    {Name:"CRUD",Syntax:"@CRUD",Category:"Persistence, repository, transactions and search",Purpose:"Opt-in GORM basic CRUD only; no guessed aggregate persistence",Mode:"GEN",Status:contract.Optional,Roadmap:"P4",Target:contract.InterfaceTarget},
    {Name:"Query",Syntax:"@Query",Category:"Persistence, repository, transactions and search",Purpose:"Explicit typed query declaration; parameterized only",Mode:"GEN/CHECK",Status:contract.Optional,Roadmap:"P4",Target:contract.MethodTarget},
    {Name:"FindBy",Syntax:"@FindBy",Category:"Persistence, repository, transactions and search",Purpose:"Allowlisted trivial finder generation; fail on ambiguity",Mode:"GEN",Status:contract.Optional,Roadmap:"P4",Target:contract.MethodTarget},
    {Name:"Transactional",Syntax:"@Transactional",Category:"Persistence, repository, transactions and search",Purpose:"Requires explicit core/pipeline/service.Transactional wrapping for transactional behavior",Mode:"CHECK",Status:contract.Pilot,Roadmap:"P4",Target:contract.MethodTarget},
    {Name:"ReadOnly",Syntax:"@ReadOnly",Category:"Persistence, repository, transactions and search",Purpose:"Read-only transaction intent; enforce at DB where possible",Mode:"WRAP",Status:contract.Optional,Roadmap:"P4",Target:contract.MethodTarget},
    {Name:"OptimisticLock",Syntax:"@OptimisticLock",Category:"Persistence, repository, transactions and search",Purpose:"Version predicate plus conflict mapping",Mode:"GEN",Status:contract.Optional,Roadmap:"P4",Target:contract.MethodTarget},
    {Name:"TenantScoped",Syntax:"@TenantScoped",Category:"Persistence, repository, transactions and search",Purpose:"Require trusted tenant context in repository predicates",Mode:"CHECK/WRAP",Status:contract.Proposed,Roadmap:"P7",Target:contract.InterfaceTarget},
    {Name:"SoftDelete",Syntax:"@SoftDelete",Category:"Persistence, repository, transactions and search",Purpose:"Optional GORM soft-delete integration",Mode:"GEN",Status:contract.Optional,Roadmap:"P4",Target:contract.InterfaceTarget},
    {Name:"Searchable",Syntax:"@Searchable impl=Repository entity=Type config=searchConfig alias=p",Category:"Persistence, repository, transactions and search",Purpose:"Generate typed search helper for an annotated concrete repository Search method, backed by an explicit core/search configuration",Mode:"GEN",Status:contract.Pilot,Roadmap:"P4",Target:contract.MethodTarget},
    {Name:"Filterable",Syntax:"@Filterable fields=field1,field2",Category:"Persistence, repository, transactions and search",Purpose:"Check configured allowed filter field names; operator policy remains in validators",Mode:"CHECK",Status:contract.Pilot,Roadmap:"P4",Target:contract.MethodTarget},
    {Name:"Sortable",Syntax:"@Sortable fields=field1,field2",Category:"Persistence, repository, transactions and search",Purpose:"Check configured allowed sorting fields",Mode:"CHECK",Status:contract.Pilot,Roadmap:"P4",Target:contract.MethodTarget},
    {Name:"Cursor",Syntax:"@Cursor field=gid",Category:"Persistence, repository, transactions and search",Purpose:"Check the stable cursor field matches the configured sort cursor",Mode:"CHECK",Status:contract.Pilot,Roadmap:"P4",Target:contract.MethodTarget},
    {Name:"Migration",Syntax:"@Migration",Category:"Persistence, repository, transactions and search",Purpose:"Only mark migration files for linting; never auto-apply DDL",Mode:"CHECK",Status:contract.Discouraged,Roadmap:"P4",Target:contract.MethodTarget},
}
