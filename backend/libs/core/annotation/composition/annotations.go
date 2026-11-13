// Package composition owns the annotation metadata for this backend responsibility.
package composition

import "grovehub/libs/core/annotation/contract"

var UseCases = []contract.UseCase{
    {Name:"Service",Syntax:"@Service",Category:"Services, composition and dependency injection",Purpose:"Register service contract/constructor for generated composition",Mode:"CHECK",Status:contract.Pilot,Roadmap:"P5",Target:contract.TypeTarget},
    {Name:"Component",Syntax:"@Component",Category:"Services, composition and dependency injection",Purpose:"Generic registered component when no narrower role",Mode:"REG",Status:contract.Optional,Roadmap:"P5",Target:contract.TypeTarget},
    {Name:"Inject",Syntax:"@Inject",Category:"Services, composition and dependency injection",Purpose:"Constructor parameter binding, not reflection field injection",Mode:"REG",Status:contract.Optional,Roadmap:"P5",Target:contract.TypeTarget},
    {Name:"Provides",Syntax:"@Provides",Category:"Services, composition and dependency injection",Purpose:"Explicit provider function for external dependencies",Mode:"REG",Status:contract.Proposed,Roadmap:"P5",Target:contract.FunctionTarget},
    {Name:"Named",Syntax:"@Named",Category:"Services, composition and dependency injection",Purpose:"Qualify dependencies by stable identifier",Mode:"REG",Status:contract.Optional,Roadmap:"P5",Target:contract.TypeTarget},
    {Name:"Scope",Syntax:"@Scope",Category:"Services, composition and dependency injection",Purpose:"Lifetime metadata; default service-instance scope",Mode:"REG",Status:contract.Optional,Roadmap:"P5",Target:contract.TypeTarget},
    {Name:"PostConstruct",Syntax:"@PostConstruct",Category:"Services, composition and dependency injection",Purpose:"Explicit startup hook; ordered and failure-aware",Mode:"REG",Status:contract.Optional,Roadmap:"P5",Target:contract.FunctionTarget},
    {Name:"PreDestroy",Syntax:"@PreDestroy",Category:"Services, composition and dependency injection",Purpose:"Explicit close/shutdown hook with deadlines",Mode:"REG",Status:contract.Optional,Roadmap:"P5",Target:contract.FunctionTarget},
    {Name:"FeatureFlag",Syntax:"@FeatureFlag",Category:"Services, composition and dependency injection",Purpose:"Register typed feature gate; not authorization",Mode:"WRAP/CHECK",Status:contract.Optional,Roadmap:"P7",Target:contract.TypeTarget},
    {Name:"Idempotent",Syntax:"@Idempotent",Category:"Services, composition and dependency injection",Purpose:"Apply declared key scope/result/retention rules at service boundary",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P4",Target:contract.FunctionTarget},
}
