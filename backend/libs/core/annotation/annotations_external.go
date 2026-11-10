// Phase 1 typed http clients, provider adapters and external services annotation metadata. No runtime behavior.
package annotation

var externalUseCases = []UseCase{
    {Name:"Client",Syntax:"@Client",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Typed in-process or HTTP client descriptor / adapter",Mode:"GEN/REG",Status:Pilot,Roadmap:"P3",Target:InterfaceTarget},
    {Name:"ExternalService",Syntax:"@ExternalService",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Explicit third-party boundary and provider ownership",Mode:"REG",Status:Proposed,Roadmap:"P5",Target:InterfaceTarget},
    {Name:"Provider",Syntax:"@Provider",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Provider identity, capabilities and factory registration",Mode:"REG",Status:Proposed,Roadmap:"P5",Target:InterfaceTarget},
    {Name:"Adapter",Syntax:"@Adapter",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Bind vendor implementation to stable provider interface",Mode:"REG",Status:Proposed,Roadmap:"P5",Target:InterfaceTarget},
    {Name:"Endpoint",Syntax:"@Endpoint",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Typed base URL/route metadata from configuration",Mode:"REG",Status:Proposed,Roadmap:"P5",Target:MethodTarget},
    {Name:"HTTPMethod",Syntax:"@HTTPMethod",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Outgoing request verb/path mapping",Mode:"GEN",Status:Optional,Roadmap:"P5",Target:MethodTarget},
    {Name:"OAuth2",Syntax:"@OAuth2",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"OAuth flow configuration metadata, no hidden token grant",Mode:"REG",Status:Optional,Roadmap:"P5",Target:InterfaceTarget},
    {Name:"Credential",Syntax:"@Credential",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Scoped credential reference/injection; redact",Mode:"REG",Status:Proposed,Roadmap:"P5",Target:InterfaceTarget},
    {Name:"Retry",Syntax:"@Retry",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Bounded idempotent retry policy, explicit status set",Mode:"WRAP",Status:Proposed,Roadmap:"P5",Target:MethodTarget},
    {Name:"Backoff",Syntax:"@Backoff",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Deterministic backoff with jitter/cap",Mode:"WRAP",Status:Proposed,Roadmap:"P5",Target:MethodTarget},
    {Name:"CircuitBreaker",Syntax:"@CircuitBreaker",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Bounded failure protection per provider/tenant",Mode:"WRAP",Status:Optional,Roadmap:"P5",Target:MethodTarget},
    {Name:"Bulkhead",Syntax:"@Bulkhead",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Limit concurrent provider calls",Mode:"WRAP",Status:Optional,Roadmap:"P5",Target:MethodTarget},
    {Name:"Fallback",Syntax:"@Fallback",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Typed fallback function for safe error types",Mode:"WRAP",Status:Optional,Roadmap:"P5",Target:MethodTarget},
    {Name:"Webhook",Syntax:"@Webhook",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Webhook route/handler registration metadata",Mode:"REG",Status:Proposed,Roadmap:"P5",Target:InterfaceTarget},
    {Name:"VerifySignature",Syntax:"@VerifySignature",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Invoke provider-specific authenticity verifier",Mode:"WRAP",Status:Proposed,Roadmap:"P5",Target:MethodTarget},
    {Name:"HealthProbe",Syntax:"@HealthProbe",Category:"Typed HTTP clients, provider adapters and external services",Purpose:"Explicit provider health check with timeout",Mode:"REG",Status:Optional,Roadmap:"P5",Target:MethodTarget},
}
