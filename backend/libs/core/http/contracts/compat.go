// Package contracts is a compatibility facade.
// Deprecated: use the root core/http package for new code.
package contracts

import (
 corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http"
 coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
)
type Success[T any] = corehttp.Success[T]
type Error[M any] = corehttp.Error[M]
func NewSuccess[T any](data T) Success[T] { return corehttp.NewSuccess(data) }
func NewError[M any](code coreerrors.Code, message string, status int, metadata M) Error[M] { return corehttp.NewError(code,message,status,metadata) }
