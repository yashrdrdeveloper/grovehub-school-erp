// Package server owns the annotation metadata for this backend responsibility.
package server

import "grovehub/libs/core/annotation/contract"

var UseCases = []contract.UseCase{
    {Name:"Module",Syntax:"@Module",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Declare module composition entry, dependencies and exports",Mode:"REG",Status:contract.Proposed,Roadmap:"P3",Target:contract.FunctionTarget},
    {Name:"Server",Syntax:"@Server",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"HTTP server policy (bind address, limits, TLS)",Mode:"REG/CHECK",Status:contract.Optional,Roadmap:"P3",Target:contract.FunctionTarget},
    {Name:"HealthCheck",Syntax:"@HealthCheck",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Register liveness probe function",Mode:"REG",Status:contract.Proposed,Roadmap:"P3",Target:contract.FunctionTarget},
    {Name:"ReadyCheck",Syntax:"@ReadyCheck",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Register dependencies/readiness policy",Mode:"REG",Status:contract.Proposed,Roadmap:"P3",Target:contract.FunctionTarget},
    {Name:"Startup",Syntax:"@Startup",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Register ordered initializer and failure behavior",Mode:"REG",Status:contract.Optional,Roadmap:"P3",Target:contract.FunctionTarget},
    {Name:"GracefulShutdown",Syntax:"@GracefulShutdown",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Stop accepting requests, drain work, close dependencies",Mode:"REG",Status:contract.Proposed,Roadmap:"P3",Target:contract.FunctionTarget},
    {Name:"Shutdown",Syntax:"@Shutdown",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Module close handler and deadline",Mode:"REG",Status:contract.Proposed,Roadmap:"P3",Target:contract.FunctionTarget},
    {Name:"Scheduled",Syntax:"@Scheduled",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Schedule registered function with timezone/cadence",Mode:"REG",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"Cron",Syntax:"@Cron",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Cron expression with explicit timezone/misfire handling",Mode:"REG",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"LeaderOnly",Syntax:"@LeaderOnly",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Distributed lock/lease requirement for singleton jobs",Mode:"WRAP",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"BackgroundJob",Syntax:"@BackgroundJob",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Registered worker managed by lifecycle, not hidden goroutine",Mode:"REG",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
}
