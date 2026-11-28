// Package response is a compatibility facade.
// Deprecated: use the root core/http package for new code.
package response

import ("context"; "net/http"; corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http"; coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors")
func WriteSuccess[T any](writer http.ResponseWriter,status int,data T)error{return corehttp.WriteSuccess(writer,status,data)}
func WriteError(writer http.ResponseWriter,status int,definition coreerrors.Definition)error{return corehttp.WriteError(writer,status,definition)}
func WriteErrorContext(ctx context.Context,writer http.ResponseWriter,status int,definition coreerrors.Definition)error{return corehttp.WriteErrorContext(ctx,writer,status,definition)}
func WriteErrorWithMetadata[M any](writer http.ResponseWriter,status int,definition coreerrors.Definition,metadata M)error{return corehttp.WriteErrorWithMetadata(writer,status,definition,metadata)}
func WriteErrorWithMetadataContext[M any](ctx context.Context,writer http.ResponseWriter,status int,definition coreerrors.Definition,metadata M)error{return corehttp.WriteErrorWithMetadataContext(ctx,writer,status,definition,metadata)}
