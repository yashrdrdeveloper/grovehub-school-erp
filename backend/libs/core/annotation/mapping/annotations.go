// mapping owns the annotation definitions for this backend responsibility.
package mapping

import "grovehub/libs/core/annotation/contract"

var UseCases = []contract.UseCase{
    {Name:"Mapper",Syntax:"@Mapper",Category:"Mappers and representation conversion",Purpose:"Generate typed mapping implementation for interface",Mode:"GEN",Status:contract.Pilot,Roadmap:"P2",Target:contract.InterfaceTarget},
    {Name:"Mapping",Syntax:"@Mapping",Category:"Mappers and representation conversion",Purpose:"Explicit source-to-target field relation / rename",Mode:"GEN",Status:contract.Pilot,Roadmap:"P2",Target:contract.MethodTarget},
    {Name:"MapIgnore",Syntax:"@MapIgnore",Category:"Mappers and representation conversion",Purpose:"Prevent source/destination field copy; required for secrets",Mode:"CHECK/GEN",Status:contract.Proposed,Roadmap:"P2",Target:contract.MethodTarget},
    {Name:"ConvertWith",Syntax:"@ConvertWith",Category:"Mappers and representation conversion",Purpose:"Call declared pure conversion function",Mode:"GEN",Status:contract.Proposed,Roadmap:"P2",Target:contract.MethodTarget},
    {Name:"EnumMapping",Syntax:"@EnumMapping",Category:"Mappers and representation conversion",Purpose:"Enum/name conversion with exhaustive unknown handling",Mode:"GEN",Status:contract.Proposed,Roadmap:"P2",Target:contract.MethodTarget},
    {Name:"AfterMapping",Syntax:"@AfterMapping",Category:"Mappers and representation conversion",Purpose:"Invoke explicit pure post-mapper, not business behavior",Mode:"GEN",Status:contract.Optional,Roadmap:"P2",Target:contract.MethodTarget},
    {Name:"CollectionMapping",Syntax:"@CollectionMapping",Category:"Mappers and representation conversion",Purpose:"Typed collection conversion with nil/empty policy",Mode:"GEN",Status:contract.Optional,Roadmap:"P2",Target:contract.MethodTarget},
    {Name:"Flatten",Syntax:"@Flatten",Category:"Mappers and representation conversion",Purpose:"Nested field mapping for proven DTO boundary",Mode:"GEN",Status:contract.Optional,Roadmap:"P2",Target:contract.MethodTarget},
    {Name:"MappingPolicy",Syntax:"@MappingPolicy",Category:"Mappers and representation conversion",Purpose:"Strict unmapped-field/null/overflow behavior",Mode:"CHECK",Status:contract.Proposed,Roadmap:"P2",Target:contract.MethodTarget},
    {Name:"MappingIgnore",Syntax:"@MappingIgnore",Category:"Mappers and representation conversion",Purpose:"Discouraged legacy spelling; canonical directive is @MapIgnore, with no separate generator behavior",Mode:"CHECK",Status:contract.Discouraged,Roadmap:"P2",Target:contract.MethodTarget},
}
