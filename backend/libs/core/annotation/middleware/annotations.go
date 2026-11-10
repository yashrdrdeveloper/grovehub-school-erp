// Package middleware owns the annotation metadata for this backend responsibility.
package middleware

import "grovehub/libs/core/annotation/contract"

var UseCases = []contract.UseCase{
    {Name:"UseMiddleware",Syntax:"@UseMiddleware",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Ordered middleware bindings from approved registry",Mode:"REG",Status:contract.Proposed,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"RequestID",Syntax:"@RequestID",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Generate/propagate unique request ID",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"CorrelationID",Syntax:"@CorrelationID",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Propagate correlation ID distinct from request ID",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"Trace",Syntax:"@Trace",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"OpenTelemetry instrumentation and context propagation",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P6",Target:contract.MethodTarget},
    {Name:"Log",Syntax:"@Log",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Structured request/operation log policy; redact by default",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P6",Target:contract.MethodTarget},
    {Name:"Recover",Syntax:"@Recover",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Panic-to-safe-HTTP response, no suppressed observability",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"Timeout",Syntax:"@Timeout",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Deadline enforcement with cancellation-aware code",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"RateLimit",Syntax:"@RateLimit",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Tenant/actor-scoped limits and bounded state",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P7",Target:contract.MethodTarget},
    {Name:"CORS",Syntax:"@CORS",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Approved-origin policy, not permissive wildcard default",Mode:"WRAP",Status:contract.Optional,Roadmap:"P7",Target:contract.MethodTarget},
    {Name:"MaxBody",Syntax:"@MaxBody",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Request size limit before decoding",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"SecureHeaders",Syntax:"@SecureHeaders",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Security headers / transport policy",Mode:"WRAP",Status:contract.Optional,Roadmap:"P7",Target:contract.MethodTarget},
    {Name:"CacheControl",Syntax:"@CacheControl",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"HTTP caching semantics (private/no-store as needed)",Mode:"WRAP",Status:contract.Optional,Roadmap:"P7",Target:contract.MethodTarget},
    {Name:"CSRF",Syntax:"@CSRF",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Browser session/unsafe-request protection where applicable",Mode:"WRAP",Status:contract.Optional,Roadmap:"P7",Target:contract.MethodTarget},
    {Name:"Compress",Syntax:"@Compress",Category:"HTTP middleware and cross-cutting web concerns",Purpose:"Compression with content-type and side-channel checks",Mode:"WRAP",Status:contract.Optional,Roadmap:"P7",Target:contract.MethodTarget},
}
