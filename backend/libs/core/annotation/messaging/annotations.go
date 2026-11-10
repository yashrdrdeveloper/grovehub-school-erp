// Package messaging owns the annotation metadata for this backend responsibility.
package messaging

import "grovehub/libs/core/annotation/contract"

var UseCases = []contract.UseCase{
    {Name:"Producer",Syntax:"@Producer",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Typed publisher binding to MessagingClient, domain topic",Mode:"REG/GEN",Status:contract.Proposed,Roadmap:"P6",Target:contract.MethodTarget},
    {Name:"Consumer",Syntax:"@Consumer",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Typed consumer registration using MessagingClient",Mode:"REG",Status:contract.Proposed,Roadmap:"P6",Target:contract.MethodTarget},
    {Name:"Subscribe",Syntax:"@Subscribe",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Stable topic/subject and handler association",Mode:"REG",Status:contract.Proposed,Roadmap:"P6",Target:contract.MethodTarget},
    {Name:"Subject",Syntax:"@Subject",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Typed subject metadata; no per-event randomized subjects",Mode:"REG/CHECK",Status:contract.Proposed,Roadmap:"P6",Target:contract.TypeTarget},
    {Name:"Stream",Syntax:"@Stream",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Stream configuration owned by Infrastructure provisioning",Mode:"REG/CHECK",Status:contract.Optional,Roadmap:"P6",Target:contract.TypeTarget},
    {Name:"Durable",Syntax:"@Durable",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Stable durable consumer identifier",Mode:"REG",Status:contract.Proposed,Roadmap:"P6",Target:contract.TypeTarget},
    {Name:"QueueGroup",Syntax:"@QueueGroup",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Consumer work-sharing group metadata",Mode:"REG",Status:contract.Optional,Roadmap:"P6",Target:contract.TypeTarget},
    {Name:"Ack",Syntax:"@Ack",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Explicit ack policy and success handling",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P6",Target:contract.MethodTarget},
    {Name:"Nack",Syntax:"@Nack",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Explicit failure retry or reject behavior",Mode:"WRAP",Status:contract.Optional,Roadmap:"P6",Target:contract.MethodTarget},
    {Name:"RetryPolicy",Syntax:"@RetryPolicy",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Consumer delivery attempt/backoff constraints",Mode:"REG",Status:contract.Proposed,Roadmap:"P6",Target:contract.TypeTarget},
    {Name:"MaxDeliver",Syntax:"@MaxDeliver",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Upper bound on attempts; poison advisory handling",Mode:"REG",Status:contract.Proposed,Roadmap:"P6",Target:contract.TypeTarget},
    {Name:"DeadLetter",Syntax:"@DeadLetter",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Explicit dead-letter/recovery workflow; not assumed native DLQ",Mode:"REG",Status:contract.Optional,Roadmap:"P6",Target:contract.TypeTarget},
    {Name:"Concurrency",Syntax:"@Concurrency",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Consumer worker limits/backpressure",Mode:"REG",Status:contract.Proposed,Roadmap:"P6",Target:contract.MethodTarget},
    {Name:"Batch",Syntax:"@Batch",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Batch processing/ack semantics",Mode:"REG",Status:contract.Optional,Roadmap:"P6",Target:contract.MethodTarget},
    {Name:"IdempotentConsumer",Syntax:"@IdempotentConsumer",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Deduplication key and state-store contract",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P6",Target:contract.MethodTarget},
    {Name:"MessageTrace",Syntax:"@MessageTrace",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"W3C trace context publish/extract",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P6",Target:contract.MethodTarget},
    {Name:"MessageTimeout",Syntax:"@MessageTimeout",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Processing timeout / ack wait coordination",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P6",Target:contract.MethodTarget},
}
