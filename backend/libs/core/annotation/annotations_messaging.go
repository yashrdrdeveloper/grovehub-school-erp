// Phase 1 messaging producer, consumer, broker topology and delivery annotation metadata. No runtime behavior.
package annotation

var messagingUseCases = []UseCase{
    {Name:"Producer",Syntax:"@Producer",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Typed publisher binding to MessagingClient, domain topic",Mode:"REG/GEN",Status:Proposed,Roadmap:"P6",Target:MethodTarget},
    {Name:"Consumer",Syntax:"@Consumer",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Typed consumer registration using MessagingClient",Mode:"REG",Status:Proposed,Roadmap:"P6",Target:MethodTarget},
    {Name:"Subscribe",Syntax:"@Subscribe",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Stable topic/subject and handler association",Mode:"REG",Status:Proposed,Roadmap:"P6",Target:MethodTarget},
    {Name:"Subject",Syntax:"@Subject",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Typed subject metadata; no per-event randomized subjects",Mode:"REG/CHECK",Status:Proposed,Roadmap:"P6",Target:TypeTarget},
    {Name:"Stream",Syntax:"@Stream",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Stream configuration owned by Infrastructure provisioning",Mode:"REG/CHECK",Status:Optional,Roadmap:"P6",Target:TypeTarget},
    {Name:"Durable",Syntax:"@Durable",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Stable durable consumer identifier",Mode:"REG",Status:Proposed,Roadmap:"P6",Target:TypeTarget},
    {Name:"QueueGroup",Syntax:"@QueueGroup",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Consumer work-sharing group metadata",Mode:"REG",Status:Optional,Roadmap:"P6",Target:TypeTarget},
    {Name:"Ack",Syntax:"@Ack",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Explicit ack policy and success handling",Mode:"WRAP",Status:Proposed,Roadmap:"P6",Target:MethodTarget},
    {Name:"Nack",Syntax:"@Nack",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Explicit failure retry or reject behavior",Mode:"WRAP",Status:Optional,Roadmap:"P6",Target:MethodTarget},
    {Name:"RetryPolicy",Syntax:"@RetryPolicy",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Consumer delivery attempt/backoff constraints",Mode:"REG",Status:Proposed,Roadmap:"P6",Target:TypeTarget},
    {Name:"MaxDeliver",Syntax:"@MaxDeliver",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Upper bound on attempts; poison advisory handling",Mode:"REG",Status:Proposed,Roadmap:"P6",Target:TypeTarget},
    {Name:"DeadLetter",Syntax:"@DeadLetter",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Explicit dead-letter/recovery workflow; not assumed native DLQ",Mode:"REG",Status:Optional,Roadmap:"P6",Target:TypeTarget},
    {Name:"Concurrency",Syntax:"@Concurrency",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Consumer worker limits/backpressure",Mode:"REG",Status:Proposed,Roadmap:"P6",Target:MethodTarget},
    {Name:"Batch",Syntax:"@Batch",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Batch processing/ack semantics",Mode:"REG",Status:Optional,Roadmap:"P6",Target:MethodTarget},
    {Name:"IdempotentConsumer",Syntax:"@IdempotentConsumer",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Deduplication key and state-store contract",Mode:"WRAP",Status:Proposed,Roadmap:"P6",Target:MethodTarget},
    {Name:"MessageTrace",Syntax:"@MessageTrace",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"W3C trace context publish/extract",Mode:"WRAP",Status:Proposed,Roadmap:"P6",Target:MethodTarget},
    {Name:"MessageTimeout",Syntax:"@MessageTimeout",Category:"Messaging producer, consumer, broker topology and delivery",Purpose:"Processing timeout / ack wait coordination",Mode:"WRAP",Status:Proposed,Roadmap:"P6",Target:MethodTarget},
}
