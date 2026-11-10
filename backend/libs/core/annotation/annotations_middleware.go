// Phase 1 http middleware and cross-cutting web concerns annotation metadata. No runtime behavior.
package annotation

var middlewareUseCases = []UseCase{
    {Name:"UseMiddleware",Syntax:"@UseMiddleware",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Ordered middleware bindings from approved registry",Mode:"REG",Status:Proposed,Roadmap:"P3",Target:MethodTarget},
    {Name:"RequestID",Syntax:"@RequestID",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Generate/propagate unique request ID",Mode:"WRAP",Status:Proposed,Roadmap:"P3",Target:MethodTarget},
    {Name:"CorrelationID",Syntax:"@CorrelationID",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Propagate correlation ID distinct from request ID",Mode:"WRAP",Status:Proposed,Roadmap:"P3",Target:MethodTarget},
    {Name:"Trace",Syntax:"@Trace",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"OpenTelemetry instrumentation and context propagation",Mode:"WRAP",Status:Proposed,Roadmap:"P6",Target:MethodTarget},
    {Name:"Log",Syntax:"@Log",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Structured request/operation log policy; redact by default",Mode:"WRAP",Status:Proposed,Roadmap:"P6",Target:MethodTarget},
    {Name:"Recover",Syntax:"@Recover",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Panic-to-safe-HTTP response, no suppressed observability",Mode:"WRAP",Status:Proposed,Roadmap:"P3",Target:MethodTarget},
    {Name:"Timeout",Syntax:"@Timeout",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Deadline enforcement with cancellation-aware code",Mode:"WRAP",Status:Proposed,Roadmap:"P3",Target:MethodTarget},
    {Name:"RateLimit",Syntax:"@RateLimit",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Tenant/actor-scoped limits and bounded state",Mode:"WRAP",Status:Proposed,Roadmap:"P7",Target:MethodTarget},
    {Name:"CORS",Syntax:"@CORS",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Approved-origin policy, not permissive wildcard default",Mode:"WRAP",Status:Optional,Roadmap:"P7",Target:MethodTarget},
    {Name:"MaxBody",Syntax:"@MaxBody",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Request size limit before decoding",Mode:"WRAP",Status:Proposed,Roadmap:"P3",Target:MethodTarget},
    {Name:"SecureHeaders",Syntax:"@SecureHeaders",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Security headers / transport policy",Mode:"WRAP",Status:Optional,Roadmap:"P7",Target:MethodTarget},
    {Name:"CacheControl",Syntax:"@CacheControl",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"HTTP caching semantics (private/no-store as needed)",Mode:"WRAP",Status:Optional,Roadmap:"P7",Target:MethodTarget},
    {Name:"CSRF",Syntax:"@CSRF",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Browser session/unsafe-request protection where applicable",Mode:"WRAP",Status:Optional,Roadmap:"P7",Target:MethodTarget},
    {Name:"Compress",Syntax:"@Compress",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Compression with content-type and side-channel checks",Mode:"WRAP",Status:Optional,Roadmap:"P7",Target:MethodTarget},
}
