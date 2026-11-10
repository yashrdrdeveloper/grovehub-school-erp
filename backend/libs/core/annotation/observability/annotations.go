// Package observability owns the annotation metadata for this backend responsibility.
package observability

import "grovehub/libs/core/annotation/contract"

var UseCases = []contract.UseCase{
    {Name:"Log",Syntax:"@Log",Category:"Logging, tracing, metrics, auditing and diagnostics",Purpose:"Structured log wrapper through `slog`/core logging",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P7",Target:contract.MethodTarget},
    {Name:"Trace",Syntax:"@Trace",Category:"Logging, tracing, metrics, auditing and diagnostics",Purpose:"OpenTelemetry span boundary with context",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P7",Target:contract.MethodTarget},
    {Name:"Span",Syntax:"@Span",Category:"Logging, tracing, metrics, auditing and diagnostics",Purpose:"Named child span with attributes allowlist",Mode:"WRAP",Status:contract.Optional,Roadmap:"P7",Target:contract.MethodTarget},
    {Name:"Timed",Syntax:"@Timed",Category:"Logging, tracing, metrics, auditing and diagnostics",Purpose:"Histogram for operation duration with bounded labels",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P7",Target:contract.MethodTarget},
    {Name:"Counted",Syntax:"@Counted",Category:"Logging, tracing, metrics, auditing and diagnostics",Purpose:"Counter for stable operation/result labels",Mode:"WRAP",Status:contract.Optional,Roadmap:"P7",Target:contract.MethodTarget},
    {Name:"Metric",Syntax:"@Metric",Category:"Logging, tracing, metrics, auditing and diagnostics",Purpose:"Register stable metric identity and unit",Mode:"REG",Status:contract.Optional,Roadmap:"P7",Target:contract.FunctionTarget},
    {Name:"Audit",Syntax:"@Audit",Category:"Logging, tracing, metrics, auditing and diagnostics",Purpose:"Security/business audit without secret exposure",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P7",Target:contract.MethodTarget},
    {Name:"Redact",Syntax:"@Redact",Category:"Logging, tracing, metrics, auditing and diagnostics",Purpose:"Hide fields from logs, traces and errors",Mode:"CHECK/WRAP",Status:contract.Proposed,Roadmap:"P7",Target:contract.MethodTarget},
    {Name:"NoLog",Syntax:"@NoLog",Category:"Logging, tracing, metrics, auditing and diagnostics",Purpose:"Suppress payload logging, preserve audit/error metadata",Mode:"CHECK",Status:contract.Proposed,Roadmap:"P7",Target:contract.MethodTarget},
    {Name:"SLO",Syntax:"@SLO",Category:"Logging, tracing, metrics, auditing and diagnostics",Purpose:"Attach objective metadata to monitored endpoint/operation",Mode:"REG",Status:contract.Optional,Roadmap:"P7",Target:contract.FunctionTarget},
    {Name:"CorrelationID",Syntax:"@CorrelationID",Category:"Logging, tracing, metrics, auditing and diagnostics",Purpose:"Carry correlation field and trace metadata",Mode:"WRAP",Status:contract.Proposed,Roadmap:"P7",Target:contract.MethodTarget},
}
