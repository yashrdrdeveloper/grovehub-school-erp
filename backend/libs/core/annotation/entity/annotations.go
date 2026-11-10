// entity owns the annotation definitions for this backend responsibility.
package entity

import "grovehub/libs/core/annotation/contract"

var UseCases = []contract.UseCase{
    {Name:"Entity",Syntax:"@Entity schema.table [mode=table]",Category:"Entity, DTO and model lifecycle",Purpose:"Generate GORM TableName only; entity structs remain passive and never receive implicit builders, getters, or setters",Mode:"GEN",Status:contract.Existing,Roadmap:"P0",Target:contract.TypeTarget},
    {Name:"Builder",Syntax:"@Builder",Category:"Entity, DTO and model lifecycle",Purpose:"Fluent typed builder; preserve field types and optional semantics",Mode:"GEN",Status:contract.Existing,Roadmap:"P0",Target:contract.TypeTarget},
    {Name:"Getter",Syntax:"@Getter",Category:"Entity, DTO and model lifecycle",Purpose:"Typed getter; suppress conflicting custom methods",Mode:"GEN",Status:contract.Existing,Roadmap:"P0",Target:contract.TypeTarget},
    {Name:"Setter",Syntax:"@Setter",Category:"Entity, DTO and model lifecycle",Purpose:"Mutable typed setter; never default for sensitive fields",Mode:"GEN",Status:contract.Existing,Roadmap:"P0",Target:contract.TypeTarget},
    {Name:"AllArgsConstructor",Syntax:"@AllArgsConstructor",Category:"Entity, DTO and model lifecycle",Purpose:"All-fields constructor when appropriate",Mode:"GEN",Status:contract.Gombok,Roadmap:"P0",Target:contract.TypeTarget},
    {Name:"NoArgsConstructor",Syntax:"@NoArgsConstructor",Category:"Entity, DTO and model lifecycle",Purpose:"Zero-argument constructor; use sparingly in Go",Mode:"GEN",Status:contract.Gombok,Roadmap:"P0",Target:contract.TypeTarget},
    {Name:"RequiredArgsConstructor",Syntax:"@RequiredArgsConstructor",Category:"Entity, DTO and model lifecycle",Purpose:"Required field constructor; evaluate Go-specific meaning",Mode:"GEN",Status:contract.Gombok,Roadmap:"P0",Target:contract.TypeTarget},
    {Name:"ToString",Syntax:"@ToString",Category:"Entity, DTO and model lifecycle",Purpose:"String formatting; redact sensitive fields",Mode:"GEN",Status:contract.Gombok,Roadmap:"P0",Target:contract.TypeTarget},
    {Name:"Equals",Syntax:"@Equals",Category:"Entity, DTO and model lifecycle",Purpose:"Equality; explicit semantics for pointers/slices/PII",Mode:"GEN",Status:contract.Gombok,Roadmap:"P0",Target:contract.TypeTarget},
    {Name:"Data",Syntax:"@Data",Category:"Entity, DTO and model lifecycle",Purpose:"Bundle DTO convenience only, NOT blanket entity setter generation",Mode:"GEN",Status:contract.Proposed,Roadmap:"P2",Target:contract.TypeTarget},
    {Name:"Immutable",Syntax:"@Immutable",Category:"Entity, DTO and model lifecycle",Purpose:"Exclude mutators; enforce construction-only fields",Mode:"CHECK/GEN",Status:contract.Proposed,Roadmap:"P2",Target:contract.TypeTarget},
    {Name:"Ignore",Syntax:"@Ignore",Category:"Entity, DTO and model lifecycle",Purpose:"Exclude selected fields from a specified generator",Mode:"CHECK",Status:contract.Proposed,Roadmap:"P2",Target:contract.FieldTarget},
    {Name:"Default",Syntax:"@Default",Category:"Entity, DTO and model lifecycle",Purpose:"Typed default on builder/config input; avoid hidden persistence defaults",Mode:"GEN",Status:contract.Optional,Roadmap:"P2",Target:contract.FieldTarget},
    {Name:"Sensitive",Syntax:"@Sensitive",Category:"Entity, DTO and model lifecycle",Purpose:"Mark secret/PII for logging and export policy checks",Mode:"CHECK",Status:contract.Proposed,Roadmap:"P7",Target:contract.FieldTarget},
    {Name:"Redact",Syntax:"@Redact",Category:"Entity, DTO and model lifecycle",Purpose:"Generate safe string/log representation or redaction metadata",Mode:"GEN/CHECK",Status:contract.Proposed,Roadmap:"P7",Target:contract.FieldTarget},
    {Name:"Version",Syntax:"@Version",Category:"Entity, DTO and model lifecycle",Purpose:"Optional mapping to explicit optimistic-lock field",Mode:"CHECK",Status:contract.Optional,Roadmap:"P4",Target:contract.FieldTarget},
    {Name:"SoftDelete",Syntax:"@SoftDelete",Category:"Entity, DTO and model lifecycle",Purpose:"Optional persistence policy; only where schema supports it",Mode:"CHECK",Status:contract.Optional,Roadmap:"P4",Target:contract.FieldTarget},
}
