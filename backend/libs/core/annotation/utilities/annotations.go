// Package utilities owns the annotation metadata for this backend responsibility.
package utilities

import "grovehub/libs/core/annotation/contract"

var UseCases = []contract.UseCase{
    {Name:"Pure",Syntax:"@Pure",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Static advisory for side-effect-free functions; lint where possible",Mode:"CHECK",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"Deprecated",Syntax:"@Deprecated",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Expose deprecation diagnostics across call sites",Mode:"CHECK",Status:contract.Proposed,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"Experimental",Syntax:"@Experimental",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Mark unstable API and compatibility expectations",Mode:"CHECK",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"Memoize",Syntax:"@Memoize",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Cache only proven pure, bounded and safe functions",Mode:"WRAP",Status:contract.Discouraged,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"Cacheable",Syntax:"@Cacheable",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Explicit cache keys, tenant partition, TTL, invalidation",Mode:"WRAP",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"CacheEvict",Syntax:"@CacheEvict",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Invalidate declared cache keys after successful mutation",Mode:"WRAP",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"CacheKey",Syntax:"@CacheKey",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Stable typed key fields, no secrets",Mode:"GEN/CHECK",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"TTL",Syntax:"@TTL",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Retention/expiry declaration with time unit",Mode:"REG",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"ConcurrencySafe",Syntax:"@ConcurrencySafe",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Document verification obligation; not an automatic race fix",Mode:"CHECK",Status:contract.Discouraged,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"Bulkhead",Syntax:"@Bulkhead",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Concurrency semaphore for external work",Mode:"WRAP",Status:contract.Optional,Roadmap:"P5",Target:contract.FunctionTarget},
    {Name:"Checksum",Syntax:"@Checksum",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Explicit integrity check for files/messages",Mode:"WRAP",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"Retention",Syntax:"@Retention",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Data lifecycle metadata; requires actual cleanup process",Mode:"REG/CHECK",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"UploadLimit",Syntax:"@UploadLimit",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Enforce file request/storage limit",Mode:"WRAP",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"Encrypt",Syntax:"@Encrypt",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Bounded encryption policy via core/crypto if added",Mode:"WRAP",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
    {Name:"GenerateStringer",Syntax:"@GenerateStringer",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Typed enum stringer when needed",Mode:"GEN",Status:contract.Optional,Roadmap:"P8",Target:contract.FunctionTarget},
}
