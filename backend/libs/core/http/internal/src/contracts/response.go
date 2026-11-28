package contracts

import coreerrors "grovehub/libs/core/errors"

type Success[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

func NewSuccess[T any](data T) Success[T] {
	return Success[T]{Success: true, Data: data}
}

type Error[M any] struct {
	Success       bool            `json:"success"`
	ErrorCode     coreerrors.Code `json:"errorCode"`
	ErrorMessage  string          `json:"errorMessage"`
	HTTPErrorCode int             `json:"httpErrorCode"`
	ErrorMetaData M               `json:"errorMetaData"`
}

func NewError[M any](code coreerrors.Code, message string, httpErrorCode int, metadata M) Error[M] {
	return Error[M]{
		Success:       false,
		ErrorCode:     code,
		ErrorMessage:  message,
		HTTPErrorCode: httpErrorCode,
		ErrorMetaData: metadata,
	}
}
