// Package http owns the annotation metadata for this backend responsibility.
package http

import "grovehub/libs/core/annotation/contract"

var UseCases = []contract.UseCase{
    {Name:"Controller",Syntax:"@Controller",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"Register typed HTTP controller/adapter",Mode:"REG",Status:contract.Proposed,Roadmap:"P3",Target:contract.InterfaceTarget},
    {Name:"Routes",Syntax:"@Routes /resource",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"Controller/resource path relative to service prefix and access tier",Mode:"GEN",Status:contract.Pilot,Roadmap:"P3",Target:contract.InterfaceTarget},
    {Name:"Routes",Syntax:"@Routes /service",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"Service-wide route prefix on routes.RegisterHTTP; requires categorized controller packages",Mode:"GEN",Status:contract.Pilot,Roadmap:"P3",Target:contract.FunctionTarget},
    {Name:"GET",Syntax:"@GET",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"GET route with typed path and handler contract",Mode:"GEN",Status:contract.Pilot,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"POST",Syntax:"@POST",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"POST route with typed handler contract",Mode:"GEN",Status:contract.Pilot,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"PUT",Syntax:"@PUT",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"PUT route with replacement semantics",Mode:"GEN",Status:contract.Pilot,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"PATCH",Syntax:"@PATCH",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"PATCH route preserving field presence",Mode:"GEN",Status:contract.Pilot,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"DELETE",Syntax:"@DELETE",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"DELETE route",Mode:"GEN",Status:contract.Pilot,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"PathParam",Syntax:"@PathParam",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"Typed path extraction with decoding/validation",Mode:"GEN",Status:contract.Optional,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"QueryParam",Syntax:"@QueryParam",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"Typed query extraction, limits and defaults",Mode:"GEN",Status:contract.Optional,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"HeaderParam",Syntax:"@HeaderParam",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"Typed header extraction; do not trust user/tenant IDs",Mode:"GEN",Status:contract.Optional,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"RequestBody",Syntax:"@RequestBody",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"Strict typed decoding/size bound",Mode:"GEN",Status:contract.Proposed,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"Response",Syntax:"@Response",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"Typed response serializer using shared envelope",Mode:"GEN",Status:contract.Proposed,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"Status",Syntax:"@Status",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"Declared success status and error mapping",Mode:"GEN/CHECK",Status:contract.Optional,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"Consumes",Syntax:"@Consumes",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"Request media-type constraints",Mode:"CHECK",Status:contract.Optional,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"Produces",Syntax:"@Produces",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"Response content-type constraints",Mode:"CHECK",Status:contract.Optional,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"OpenAPI",Syntax:"@OpenAPI",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"Generate or validate OpenAPI from contract source",Mode:"CHECK/GEN",Status:contract.Proposed,Roadmap:"P3",Target:contract.InterfaceTarget},
    {Name:"OperationID",Syntax:"@OperationID",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"Stable API operation identity",Mode:"CHECK",Status:contract.Optional,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"Example",Syntax:"@Example",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"API examples with contract validation",Mode:"CHECK",Status:contract.Optional,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"Deprecated",Syntax:"@Deprecated",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"Deprecation metadata and compatibility warning",Mode:"CHECK",Status:contract.Optional,Roadmap:"P3",Target:contract.MethodTarget},
    {Name:"Public",Syntax:"@Public",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"Explicit public endpoint exemption, security-reviewed",Mode:"CHECK",Status:contract.Proposed,Roadmap:"P7",Target:contract.MethodTarget},
    {Name:"Schema",Syntax:"@Schema",Category:"HTTP controllers, routes, request/response and API docs",Purpose:"Typed DTO/OpenAPI schema metadata; verify against authoritative public OpenAPI contract",Mode:"CHECK/GEN",Status:contract.Proposed,Roadmap:"P3",Target:contract.TypeTarget},
}
