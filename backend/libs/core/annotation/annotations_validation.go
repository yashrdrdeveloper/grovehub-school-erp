// Phase 1 validation, normalization and structured errors annotation metadata. No runtime behavior.
package annotation

var validationUseCases = []UseCase{
    {Name:"Validator",Syntax:"@Validator",Category:"Validation, normalization and structured errors",Purpose:"Typed validator registration / wrapper for DTO type",Mode:"GEN/REG",Status:Pilot,Roadmap:"P2",Target:TypeTarget},
    {Name:"Valid",Syntax:"@Valid",Category:"Validation, normalization and structured errors",Purpose:"Nested/request validation at explicit service/controller boundary",Mode:"WRAP/TAG",Status:Proposed,Roadmap:"P2",Target:FieldTarget},
    {Name:"NonNull",Syntax:"@NonNull",Category:"Validation, normalization and structured errors",Purpose:"Reject nil/absent nullable input; distinguish from blank/zero",Mode:"TAG",Status:Proposed,Roadmap:"P2",Target:FieldTarget},
    {Name:"NotBlank",Syntax:"@NotBlank",Category:"Validation, normalization and structured errors",Purpose:"Trim-aware non-whitespace text validation",Mode:"TAG",Status:Pilot,Roadmap:"P2",Target:FieldTarget},
    {Name:"NotEmpty",Syntax:"@NotEmpty",Category:"Validation, normalization and structured errors",Purpose:"Non-empty string/list/map",Mode:"TAG",Status:Proposed,Roadmap:"P2",Target:FieldTarget},
    {Name:"Length",Syntax:"@Length",Category:"Validation, normalization and structured errors",Purpose:"Rune-based or byte-based length with declared unit",Mode:"TAG",Status:Proposed,Roadmap:"P2",Target:FieldTarget},
    {Name:"Min",Syntax:"@Min",Category:"Validation, normalization and structured errors",Purpose:"Minimum numeric value or length, explicit semantics",Mode:"TAG",Status:Proposed,Roadmap:"P2",Target:FieldTarget},
    {Name:"Max",Syntax:"@Max",Category:"Validation, normalization and structured errors",Purpose:"Maximum numeric value or length, explicit semantics",Mode:"TAG",Status:Proposed,Roadmap:"P2",Target:FieldTarget},
    {Name:"Range",Syntax:"@Range",Category:"Validation, normalization and structured errors",Purpose:"Bounded numeric/rune length range",Mode:"TAG",Status:Proposed,Roadmap:"P2",Target:FieldTarget},
    {Name:"Email",Syntax:"@Email",Category:"Validation, normalization and structured errors",Purpose:"Email syntax; coordinate with normalized canonical value",Mode:"TAG",Status:Proposed,Roadmap:"P2",Target:FieldTarget},
    {Name:"Phone",Syntax:"@Phone",Category:"Validation, normalization and structured errors",Purpose:"Phone technical format and supported region policy",Mode:"TAG",Status:Proposed,Roadmap:"P2",Target:FieldTarget},
    {Name:"UUIDv7",Syntax:"@UUIDv7",Category:"Validation, normalization and structured errors",Purpose:"UUIDv7 requirement, not generic UUID-only",Mode:"TAG",Status:Proposed,Roadmap:"P2",Target:FieldTarget},
    {Name:"Pattern",Syntax:"@Pattern",Category:"Validation, normalization and structured errors",Purpose:"Constrained regex/pattern with bounded compile cost",Mode:"TAG",Status:Optional,Roadmap:"P2",Target:FieldTarget},
    {Name:"OneOf",Syntax:"@OneOf",Category:"Validation, normalization and structured errors",Purpose:"Allow-listed enum values and unknown handling",Mode:"TAG",Status:Proposed,Roadmap:"P2",Target:FieldTarget},
    {Name:"UniqueItems",Syntax:"@UniqueItems",Category:"Validation, normalization and structured errors",Purpose:"Duplicate detection with explicit key/normalization",Mode:"TAG",Status:Proposed,Roadmap:"P2",Target:FieldTarget},
    {Name:"Normalize",Syntax:"@Normalize",Category:"Validation, normalization and structured errors",Purpose:"Named normalization pipeline: trim, email, phone, UUID",Mode:"TAG/GEN",Status:Pilot,Roadmap:"P2",Target:FieldTarget},
    {Name:"RequiredIf",Syntax:"@RequiredIf",Category:"Validation, normalization and structured errors",Purpose:"Conditional field requiredness",Mode:"TAG",Status:Optional,Roadmap:"P2",Target:FieldTarget},
    {Name:"MutuallyExclusive",Syntax:"@MutuallyExclusive",Category:"Validation, normalization and structured errors",Purpose:"Cross-field XOR invariant",Mode:"TAG",Status:Optional,Roadmap:"P2",Target:FieldTarget},
    {Name:"CrossField",Syntax:"@CrossField",Category:"Validation, normalization and structured errors",Purpose:"Register custom typed cross-field rule",Mode:"REG",Status:Proposed,Roadmap:"P2",Target:TypeTarget},
    {Name:"ValidatePatch",Syntax:"@ValidatePatch",Category:"Validation, normalization and structured errors",Purpose:"Presence-aware PATCH normalization/validation",Mode:"GEN",Status:Proposed,Roadmap:"P2",Target:TypeTarget},
    {Name:"Error",Syntax:"@Error",Category:"Validation, normalization and structured errors",Purpose:"Resolve domain error definition from constants/error; not stringly typed",Mode:"CHECK/WRAP",Status:Proposed,Roadmap:"P2",Target:TypeTarget},
    {Name:"StrictJSON",Syntax:"@StrictJSON",Category:"Validation, normalization and structured errors",Purpose:"Reject unknown properties at transport boundary",Mode:"WRAP",Status:Optional,Roadmap:"P3",Target:TypeTarget},
    {Name:"UUID",Syntax:"@UUID",Category:"Validation, normalization and structured errors",Purpose:"Validate UUID syntax when any version is permitted; use @UUIDv7 when the domain requires v7",Mode:"TAG",Status:Optional,Roadmap:"P2",Target:FieldTarget},
}
