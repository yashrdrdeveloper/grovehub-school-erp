package response

import coreerrors "grovehub/libs/core/errors"

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
