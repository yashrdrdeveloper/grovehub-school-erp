// Package event owns the annotation metadata for this backend responsibility.
package event

import "grovehub/libs/core/annotation/contract"

var UseCases = []contract.UseCase{
    {Name:"Event",Syntax:"@Event",Category:"Events, envelopes and event lifecycle",Purpose:"Domain-owned topic, version and payload metadata",Mode:"GEN/CHECK",Status:contract.Pilot,Roadmap:"P6",Target:contract.TypeTarget},
    {Name:"EventVersion",Syntax:"@EventVersion",Category:"Events, envelopes and event lifecycle",Purpose:"Stable event schema version and compatibility checks",Mode:"CHECK",Status:contract.Proposed,Roadmap:"P6",Target:contract.TypeTarget},
    {Name:"EventHandler",Syntax:"@EventHandler",Category:"Events, envelopes and event lifecycle",Purpose:"Typed event handler registration",Mode:"REG",Status:contract.Proposed,Roadmap:"P6",Target:contract.MethodTarget},
    {Name:"PublishEvent",Syntax:"@PublishEvent",Category:"Events, envelopes and event lifecycle",Purpose:"Explicit event publisher decorator metadata only",Mode:"REG/WRAP",Status:contract.Optional,Roadmap:"P6",Target:contract.MethodTarget},
    {Name:"AfterCommit",Syntax:"@AfterCommit",Category:"Events, envelopes and event lifecycle",Purpose:"Bind to existing after-commit callback semantics",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P6",Target:contract.MethodTarget},
    {Name:"Outbox",Syntax:"@Outbox",Category:"Events, envelopes and event lifecycle",Purpose:"Opt-in durable outbox persistence/relay contract",Mode:"REG",Status:contract.Optional,Roadmap:"P6",Target:contract.TypeTarget},
    {Name:"EventSchema",Syntax:"@EventSchema",Category:"Events, envelopes and event lifecycle",Purpose:"Contract/schema compatibility checks",Mode:"CHECK",Status:contract.Optional,Roadmap:"P6",Target:contract.TypeTarget},
    {Name:"EventSource",Syntax:"@EventSource",Category:"Events, envelopes and event lifecycle",Purpose:"Domain source metadata, not envelope field",Mode:"REG",Status:contract.Optional,Roadmap:"P6",Target:contract.TypeTarget},
    {Name:"Causation",Syntax:"@Causation",Category:"Events, envelopes and event lifecycle",Purpose:"Propagate existing command/parent event ID",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P6",Target:contract.MethodTarget},
    {Name:"Correlation",Syntax:"@Correlation",Category:"Events, envelopes and event lifecycle",Purpose:"Propagate existing correlation/trace context",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P6",Target:contract.MethodTarget},
    {Name:"EventPII",Syntax:"@EventPII",Category:"Events, envelopes and event lifecycle",Purpose:"Classify full-snapshot events and retention policy",Mode:"CHECK",Status:contract.Proposed,Roadmap:"P7",Target:contract.TypeTarget},
}
