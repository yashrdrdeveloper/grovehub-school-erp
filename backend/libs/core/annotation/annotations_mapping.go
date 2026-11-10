// Phase 1 mappers and representation conversion annotation metadata. No runtime behavior.
package annotation

var mappingUseCases = []UseCase{
    {Name:"Mapper",Syntax:"@Mapper",Category:"Mappers and representation conversion",Purpose:"Generate typed mapping implementation for interface",Mode:"GEN",Status:Proposed,Roadmap:"P2",Target:InterfaceTarget},
    {Name:"Mapping",Syntax:"@Mapping",Category:"Mappers and representation conversion",Purpose:"Explicit source-to-target field relation / rename",Mode:"GEN",Status:Proposed,Roadmap:"P2",Target:MethodTarget},
    {Name:"MapIgnore",Syntax:"@MapIgnore",Category:"Mappers and representation conversion",Purpose:"Prevent source/destination field copy; required for secrets",Mode:"CHECK/GEN",Status:Proposed,Roadmap:"P2",Target:MethodTarget},
    {Name:"ConvertWith",Syntax:"@ConvertWith",Category:"Mappers and representation conversion",Purpose:"Call declared pure conversion function",Mode:"GEN",Status:Proposed,Roadmap:"P2",Target:MethodTarget},
    {Name:"EnumMapping",Syntax:"@EnumMapping",Category:"Mappers and representation conversion",Purpose:"Enum/name conversion with exhaustive unknown handling",Mode:"GEN",Status:Proposed,Roadmap:"P2",Target:MethodTarget},
    {Name:"AfterMapping",Syntax:"@AfterMapping",Category:"Mappers and representation conversion",Purpose:"Invoke explicit pure post-mapper, not business behavior",Mode:"GEN",Status:Optional,Roadmap:"P2",Target:MethodTarget},
    {Name:"CollectionMapping",Syntax:"@CollectionMapping",Category:"Mappers and representation conversion",Purpose:"Typed collection conversion with nil/empty policy",Mode:"GEN",Status:Optional,Roadmap:"P2",Target:MethodTarget},
    {Name:"Flatten",Syntax:"@Flatten",Category:"Mappers and representation conversion",Purpose:"Nested field mapping for proven DTO boundary",Mode:"GEN",Status:Optional,Roadmap:"P2",Target:MethodTarget},
    {Name:"MappingPolicy",Syntax:"@MappingPolicy",Category:"Mappers and representation conversion",Purpose:"Strict unmapped-field/null/overflow behavior",Mode:"CHECK",Status:Proposed,Roadmap:"P2",Target:MethodTarget},
    {Name:"MappingIgnore",Syntax:"@MappingIgnore",Category:"Mappers and representation conversion",Purpose:"Discouraged legacy spelling; canonical directive is @MapIgnore, with no separate generator behavior",Mode:"CHECK",Status:Discouraged,Roadmap:"P2",Target:MethodTarget},
}
