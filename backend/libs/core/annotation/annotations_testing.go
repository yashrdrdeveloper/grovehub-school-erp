// Phase 1 testing, fixtures, mocks, contract tests and benchmarks annotation metadata. No runtime behavior.
package annotation

var testingUseCases = []UseCase{
    {Name:"UnitTest",Syntax:"@UnitTest",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Test category metadata for Ginkgo suite grouping",Mode:"CHECK",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"IntegrationTest",Syntax:"@IntegrationTest",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Service/database/broker integration category",Mode:"CHECK",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"E2E",Syntax:"@E2E",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Full HTTP/worker flow category",Mode:"CHECK",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"Fixture",Syntax:"@Fixture",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Typed deterministic fixture constructor",Mode:"GEN",Status:Proposed,Roadmap:"P8",Target:FunctionTarget},
    {Name:"Mock",Syntax:"@Mock",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Generate interface fake/mock without expectation magic",Mode:"GEN",Status:Proposed,Roadmap:"P8",Target:FunctionTarget},
    {Name:"GenerateMock",Syntax:"@GenerateMock",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Explicit mock target and stable generated location",Mode:"GEN",Status:Proposed,Roadmap:"P8",Target:FunctionTarget},
    {Name:"ContractTest",Syntax:"@ContractTest",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Cross-module OpenAPI/client/event contract test manifest",Mode:"REG/GEN",Status:Proposed,Roadmap:"P8",Target:FunctionTarget},
    {Name:"PropertyTest",Syntax:"@PropertyTest",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Property/fuzz test input metadata",Mode:"GEN",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"TestMatrix",Syntax:"@TestMatrix",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Explicit table-driven variations",Mode:"GEN",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"TestTimeout",Syntax:"@TestTimeout",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Timeout for test suite/fixture",Mode:"REG",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"RequiresDocker",Syntax:"@RequiresDocker",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Skip/fail rationale when containers unavailable",Mode:"CHECK",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"Benchmark",Syntax:"@Benchmark",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Register performance benchmark and baseline",Mode:"REG",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"Seed",Syntax:"@Seed",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Deterministic seed/clock and resource cleanup",Mode:"REG",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"Golden",Syntax:"@Golden",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Golden-data drift check with explicit update step",Mode:"CHECK",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
}
