// Phase 1 utilities, helpers, cache, concurrency and shared contracts annotation metadata. No runtime behavior.
package annotation

var utilitiesUseCases = []UseCase{
    {Name:"Pure",Syntax:"@Pure",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Static advisory for side-effect-free functions; lint where possible",Mode:"CHECK",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"Deprecated",Syntax:"@Deprecated",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Expose deprecation diagnostics across call sites",Mode:"CHECK",Status:Proposed,Roadmap:"P8",Target:FunctionTarget},
    {Name:"Experimental",Syntax:"@Experimental",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Mark unstable API and compatibility expectations",Mode:"CHECK",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"Memoize",Syntax:"@Memoize",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Cache only proven pure, bounded and safe functions",Mode:"WRAP",Status:Discouraged,Roadmap:"P8",Target:FunctionTarget},
    {Name:"Cacheable",Syntax:"@Cacheable",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Explicit cache keys, tenant partition, TTL, invalidation",Mode:"WRAP",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"CacheEvict",Syntax:"@CacheEvict",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Invalidate declared cache keys after successful mutation",Mode:"WRAP",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"CacheKey",Syntax:"@CacheKey",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Stable typed key fields, no secrets",Mode:"GEN/CHECK",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"TTL",Syntax:"@TTL",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Retention/expiry declaration with time unit",Mode:"REG",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"ConcurrencySafe",Syntax:"@ConcurrencySafe",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Document verification obligation; not an automatic race fix",Mode:"CHECK",Status:Discouraged,Roadmap:"P8",Target:FunctionTarget},
    {Name:"Bulkhead",Syntax:"@Bulkhead",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Concurrency semaphore for external work",Mode:"WRAP",Status:Optional,Roadmap:"P5",Target:FunctionTarget},
    {Name:"Checksum",Syntax:"@Checksum",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Explicit integrity check for files/messages",Mode:"WRAP",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"Retention",Syntax:"@Retention",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Data lifecycle metadata; requires actual cleanup process",Mode:"REG/CHECK",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"UploadLimit",Syntax:"@UploadLimit",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Enforce file request/storage limit",Mode:"WRAP",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"Encrypt",Syntax:"@Encrypt",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Bounded encryption policy via core/crypto if added",Mode:"WRAP",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
    {Name:"GenerateStringer",Syntax:"@GenerateStringer",Category:"Utilities, helpers, cache, concurrency and shared contracts",Purpose:"Typed enum stringer when needed",Mode:"GEN",Status:Optional,Roadmap:"P8",Target:FunctionTarget},
}
