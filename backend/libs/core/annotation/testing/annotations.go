// Package testing owns the annotation metadata for this backend responsibility.
package testing

import "grovehub/libs/core/annotation/contract"

var UseCases = []contract.UseCase{
    {Name:"Testing",Syntax:"@Testing suite=Name",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Generate the Ginkgo suite entry point from an annotated spec registration function",Mode:"GEN",Status:contract.Existing,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"BeforeAll",Syntax:"@BeforeAll",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Register a Ginkgo BeforeSuite fixture once per test suite",Mode:"REG",Status:contract.Existing,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"BeforeEach",Syntax:"@BeforeEach",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Register per-example Ginkgo setup inside generated suite",Mode:"REG",Status:contract.Existing,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"UnitTest",Syntax:"@UnitTest",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Test category metadata for Ginkgo suite grouping",Mode:"CHECK",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"IntegrationTest",Syntax:"@IntegrationTest",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Service/database/broker integration category",Mode:"CHECK",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"E2E",Syntax:"@E2E",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Full HTTP/worker flow category",Mode:"CHECK",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"Fixture",Syntax:"@Fixture",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Typed deterministic fixture constructor",Mode:"GEN",Status:contract.Proposed,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"Mock",Syntax:"@Mock",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Generate a typed GoMock for an explicitly declared test interface",Mode:"GEN",Status:contract.Existing,Roadmap:"P8",Target:contract.InterfaceTarget},
    {Name:"GenerateMock",Syntax:"@GenerateMock",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Explicit mock target and stable generated location",Mode:"GEN",Status:contract.Proposed,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"ContractTest",Syntax:"@ContractTest",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Cross-module OpenAPI/client/event contract test manifest",Mode:"REG/GEN",Status:contract.Proposed,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"PropertyTest",Syntax:"@PropertyTest",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Property/fuzz test input metadata",Mode:"GEN",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"TestMatrix",Syntax:"@TestMatrix",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Explicit table-driven variations",Mode:"GEN",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"TestTimeout",Syntax:"@TestTimeout",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Timeout for test suite/fixture",Mode:"REG",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"RequiresDocker",Syntax:"@RequiresDocker",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Skip/fail rationale when containers unavailable",Mode:"CHECK",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"Benchmark",Syntax:"@Benchmark",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Register performance benchmark and baseline",Mode:"REG",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"Seed",Syntax:"@Seed",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Deterministic seed/clock and resource cleanup",Mode:"REG",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"Golden",Syntax:"@Golden",Category:"Testing, fixtures, mocks, contract tests and benchmarks",Purpose:"Golden-data drift check with explicit update step",Mode:"CHECK",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
}
