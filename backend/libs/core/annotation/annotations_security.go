// Phase 1 security, secrets, tenancy and policy annotation metadata. No runtime behavior.
package annotation

var securityUseCases = []UseCase{
    {Name:"Authenticate",Syntax:"@Authenticate",Category:"Security, secrets, tenancy and policy",Purpose:"Require verified identity; never trust caller-filled principal",Mode:"WRAP/CHECK",Status:Proposed,Roadmap:"P7",Target:MethodTarget},
    {Name:"Authorize",Syntax:"@Authorize",Category:"Security, secrets, tenancy and policy",Purpose:"Evaluate business permission against trusted actor/context",Mode:"WRAP",Status:Proposed,Roadmap:"P7",Target:MethodTarget},
    {Name:"RolesAllowed",Syntax:"@RolesAllowed",Category:"Security, secrets, tenancy and policy",Purpose:"Role-based rule using central IAM policy",Mode:"WRAP",Status:Optional,Roadmap:"P7",Target:MethodTarget},
    {Name:"RequiresPermission",Syntax:"@RequiresPermission",Category:"Security, secrets, tenancy and policy",Purpose:"Named policy requirement with resource context",Mode:"WRAP",Status:Proposed,Roadmap:"P7",Target:MethodTarget},
    {Name:"OwnerOnly",Syntax:"@OwnerOnly",Category:"Security, secrets, tenancy and policy",Purpose:"Ownership check must resolve resource via service safely",Mode:"WRAP",Status:Optional,Roadmap:"P7",Target:MethodTarget},
    {Name:"TenantScoped",Syntax:"@TenantScoped",Category:"Security, secrets, tenancy and policy",Purpose:"Enforce tenant context across route/service/repository",Mode:"CHECK/WRAP",Status:Proposed,Roadmap:"P7",Target:MethodTarget},
    {Name:"Sensitive",Syntax:"@Sensitive",Category:"Security, secrets, tenancy and policy",Purpose:"Data classification for logs/contracts/events",Mode:"CHECK",Status:Proposed,Roadmap:"P7",Target:FieldTarget},
    {Name:"Redact",Syntax:"@Redact",Category:"Security, secrets, tenancy and policy",Purpose:"Redaction of diagnostics/log output",Mode:"WRAP/GEN",Status:Proposed,Roadmap:"P7",Target:FieldTarget},
    {Name:"Secret",Syntax:"@Secret",Category:"Security, secrets, tenancy and policy",Purpose:"Config/credential never exposed in diagnostics",Mode:"CHECK",Status:Proposed,Roadmap:"P5",Target:FieldTarget},
    {Name:"Encrypt",Syntax:"@Encrypt",Category:"Security, secrets, tenancy and policy",Purpose:"Encryption policy for stored fields, provider-owned crypto",Mode:"WRAP",Status:Optional,Roadmap:"P7",Target:FieldTarget},
    {Name:"Audit",Syntax:"@Audit",Category:"Security, secrets, tenancy and policy",Purpose:"Emit bounded, access-controlled audit record",Mode:"WRAP",Status:Proposed,Roadmap:"P6",Target:MethodTarget},
}
