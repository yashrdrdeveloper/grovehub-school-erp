// Package endpoint is a compatibility facade.
// Deprecated: use the root core/http package for new code.
package endpoint

import ("context"; "net/http"; corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http"; corerequest "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/request")
type Request[Path,Query,Body any] = corehttp.Request[Path,Query,Body]
func Call[T any](ctx context.Context, handler http.Handler, method,route string, pathDTO,body any)(T,int,error){return corehttp.Call[T](ctx,handler,method,route,pathDTO,body)}
func CallRequest[Response,Path,Query,Body any](ctx context.Context,handler http.Handler,method,route string,input Request[Path,Query,Body])(Response,int,error){return corehttp.CallRequest[Response](ctx,handler,method,route,input)}
func CallStruct[Response any, Input corerequest.Contract](ctx context.Context,handler http.Handler,method,route string,input Input)(Response,int,error){return corehttp.CallStruct[Response](ctx,handler,method,route,input)}
