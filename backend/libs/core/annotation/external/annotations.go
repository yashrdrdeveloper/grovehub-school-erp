// Package external owns the annotation metadata for this backend responsibility.
package external

import "grovehub/libs/core/annotation/contract"

var UseCases = []contract.UseCase{
    {Name:"Client",Syntax:"@Client",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Typed in-process or HTTP client descriptor / adapter",Mode:"GEN/REG",Status:contract.Pilot,Roadmap:"P3",Target:contract.InterfaceTarget},
    {Name:"ExternalService",Syntax:"@ExternalService",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Explicit third-party boundary and provider ownership",Mode:"REG",Status:contract.Pilot,Roadmap:"P5",Target:contract.InterfaceTarget},
    {Name:"Provider",Syntax:"@Provider",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Provider identity, capabilities and factory registration",Mode:"REG",Status:contract.Proposed,Roadmap:"P5",Target:contract.InterfaceTarget},
    {Name:"Adapter",Syntax:"@Adapter",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Bind vendor implementation to stable provider interface",Mode:"REG",Status:contract.Proposed,Roadmap:"P5",Target:contract.InterfaceTarget},
    {Name:"Endpoint",Syntax:"@Endpoint",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Typed base URL/route metadata from configuration",Mode:"REG",Status:contract.Proposed,Roadmap:"P5",Target:contract.MethodTarget},
    {Name:"HTTPMethod",Syntax:"@HTTPMethod",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Outgoing request verb/path mapping",Mode:"GEN",Status:contract.Optional,Roadmap:"P5",Target:contract.MethodTarget},
    {Name:"OAuth2",Syntax:"@OAuth2",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"OAuth flow configuration metadata, no hidden token grant",Mode:"REG",Status:contract.Optional,Roadmap:"P5",Target:contract.InterfaceTarget},
    {Name:"Credential",Syntax:"@Credential",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Scoped credential reference/injection; redact",Mode:"REG",Status:contract.Proposed,Roadmap:"P5",Target:contract.InterfaceTarget},
    {Name:"Retry",Syntax:"@Retry",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Bounded idempotent retry policy, explicit status set",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P5",Target:contract.MethodTarget},
    {Name:"Backoff",Syntax:"@Backoff",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Deterministic backoff with jitter/cap",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P5",Target:contract.MethodTarget},
    {Name:"CircuitBreaker",Syntax:"@CircuitBreaker",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Bounded failure protection per provider/tenant",Mode:"WRAP",Status:contract.Optional,Roadmap:"P5",Target:contract.MethodTarget},
    {Name:"Bulkhead",Syntax:"@Bulkhead",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Limit concurrent provider calls",Mode:"WRAP",Status:contract.Optional,Roadmap:"P5",Target:contract.MethodTarget},
    {Name:"Fallback",Syntax:"@Fallback",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Typed fallback function for safe error types",Mode:"WRAP",Status:contract.Optional,Roadmap:"P5",Target:contract.MethodTarget},
    {Name:"Webhook",Syntax:"@Webhook",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Webhook route/handler registration metadata",Mode:"REG",Status:contract.Proposed,Roadmap:"P5",Target:contract.InterfaceTarget},
    {Name:"VerifySignature",Syntax:"@VerifySignature",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Invoke provider-specific authenticity verifier",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P5",Target:contract.MethodTarget},
    {Name:"HealthProbe",Syntax:"@HealthProbe",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Explicit provider health check with timeout",Mode:"REG",Status:contract.Optional,Roadmap:"P5",Target:contract.MethodTarget},
}
