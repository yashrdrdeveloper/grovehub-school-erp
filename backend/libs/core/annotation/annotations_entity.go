// Phase 1 entity, dto and model lifecycle annotation metadata. No runtime behavior.
package annotation

var entityUseCases = []UseCase{
    {Name:"Entity",Syntax:"@Entity schema.table [mode=table]",Category:"Entity, DTO and model lifecycle",Purpose:"Generate API-compatible entity builders, safe accessors and GORM TableName; mode=table omits builders/accessors for sensitive models",Mode:"GEN",Status:Existing,Roadmap:"P0",Target:TypeTarget},
    {Name:"Builder",Syntax:"@Builder",Category:"Entity, DTO and model lifecycle",Purpose:"Fluent typed builder; preserve field types and optional semantics",Mode:"GEN",Status:Existing,Roadmap:"P0",Target:TypeTarget},
    {Name:"Getter",Syntax:"@Getter",Category:"Entity, DTO and model lifecycle",Purpose:"Typed getter; suppress conflicting custom methods",Mode:"GEN",Status:Existing,Roadmap:"P0",Target:TypeTarget},
    {Name:"Setter",Syntax:"@Setter",Category:"Entity, DTO and model lifecycle",Purpose:"Mutable typed setter; never default for sensitive fields",Mode:"GEN",Status:Existing,Roadmap:"P0",Target:TypeTarget},
    {Name:"AllArgsConstructor",Syntax:"@AllArgsConstructor",Category:"Entity, DTO and model lifecycle",Purpose:"All-fields constructor when appropriate",Mode:"GEN",Status:Gombok,Roadmap:"P0",Target:TypeTarget},
    {Name:"NoArgsConstructor",Syntax:"@NoArgsConstructor",Category:"Entity, DTO and model lifecycle",Purpose:"Zero-argument constructor; use sparingly in Go",Mode:"GEN",Status:Gombok,Roadmap:"P0",Target:TypeTarget},
    {Name:"RequiredArgsConstructor",Syntax:"@RequiredArgsConstructor",Category:"Entity, DTO and model lifecycle",Purpose:"Required field constructor; evaluate Go-specific meaning",Mode:"GEN",Status:Gombok,Roadmap:"P0",Target:TypeTarget},
    {Name:"ToString",Syntax:"@ToString",Category:"Entity, DTO and model lifecycle",Purpose:"String formatting; redact sensitive fields",Mode:"GEN",Status:Gombok,Roadmap:"P0",Target:TypeTarget},
    {Name:"Equals",Syntax:"@Equals",Category:"Entity, DTO and model lifecycle",Purpose:"Equality; explicit semantics for pointers/slices/PII",Mode:"GEN",Status:Gombok,Roadmap:"P0",Target:TypeTarget},
    {Name:"Data",Syntax:"@Data",Category:"Entity, DTO and model lifecycle",Purpose:"Bundle DTO convenience only, NOT blanket entity setter generation",Mode:"GEN",Status:Proposed,Roadmap:"P2",Target:TypeTarget},
    {Name:"Immutable",Syntax:"@Immutable",Category:"Entity, DTO and model lifecycle",Purpose:"Exclude mutators; enforce construction-only fields",Mode:"CHECK/GEN",Status:Proposed,Roadmap:"P2",Target:TypeTarget},
    {Name:"Ignore",Syntax:"@Ignore",Category:"Entity, DTO and model lifecycle",Purpose:"Exclude selected fields from a specified generator",Mode:"CHECK",Status:Proposed,Roadmap:"P2",Target:FieldTarget},
    {Name:"Default",Syntax:"@Default",Category:"Entity, DTO and model lifecycle",Purpose:"Typed default on builder/config input; avoid hidden persistence defaults",Mode:"GEN",Status:Optional,Roadmap:"P2",Target:FieldTarget},
    {Name:"Sensitive",Syntax:"@Sensitive",Category:"Entity, DTO and model lifecycle",Purpose:"Mark secret/PII for logging and export policy checks",Mode:"CHECK",Status:Proposed,Roadmap:"P7",Target:FieldTarget},
    {Name:"Redact",Syntax:"@Redact",Category:"Entity, DTO and model lifecycle",Purpose:"Generate safe string/log representation or redaction metadata",Mode:"GEN/CHECK",Status:Proposed,Roadmap:"P7",Target:FieldTarget},
    {Name:"Version",Syntax:"@Version",Category:"Entity, DTO and model lifecycle",Purpose:"Optional mapping to explicit optimistic-lock field",Mode:"CHECK",Status:Optional,Roadmap:"P4",Target:FieldTarget},
    {Name:"SoftDelete",Syntax:"@SoftDelete",Category:"Entity, DTO and model lifecycle",Purpose:"Optional persistence policy; only where schema supports it",Mode:"CHECK",Status:Optional,Roadmap:"P4",Target:FieldTarget},
}
