// Phase 1 server, launcher, startup, shutdown and background jobs annotation metadata. No runtime behavior.
package annotation

var serverUseCases = []UseCase{
    {Name:"Module",Syntax:"@Module",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Declare module composition entry, dependencies and exports",Mode:"REG",Status:Proposed,Roadmap:"P3",Target:FunctionTarget},
    {Name:"Server",Syntax:"@Server",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"HTTP server policy (bind address, limits, TLS)",Mode:"REG/CHECK",Status:Optional,Roadmap:"P3",Target:FunctionTarget},
    {Name:"HealthCheck",Syntax:"@HealthCheck",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Register liveness probe function",Mode:"REG",Status:Proposed,Roadmap:"P3",Target:FunctionTarget},
    {Name:"ReadyCheck",Syntax:"@ReadyCheck",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Register dependencies/readiness policy",Mode:"REG",Status:Proposed,Roadmap:"P3",Target:FunctionTarget},
    {Name:"Startup",Syntax:"@Startup",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Register ordered initializer and failure behavior",Mode:"REG",Status:Optional,Roadmap:"P3",Target:FunctionTarget},
    {Name:"GracefulShutdown",Syntax:"@GracefulShutdown",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Stop accepting requests, drain work, close dependencies",Mode:"REG",Status:Proposed,Roadmap:"P3",Target:FunctionTarget},
    {Name:"Shutdown",Syntax:"@Shutdown",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Module close handler and deadline",Mode:"REG",Status:Proposed,Roadmap:"P3",Target:FunctionTarget},
    {Name:"Scheduled",Syntax:"@Scheduled",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Schedule registered function with timezone/cadence",Mode:"REG",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"Cron",Syntax:"@Cron",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Cron expression with explicit timezone/misfire handling",Mode:"REG",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"LeaderOnly",Syntax:"@LeaderOnly",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Distributed lock/lease requirement for singleton jobs",Mode:"WRAP",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"BackgroundJob",Syntax:"@BackgroundJob",Category:"Server, launcher, startup, shutdown and background jobs",Purpose:"Registered worker managed by lifecycle, not hidden goroutine",Mode:"REG",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
}
