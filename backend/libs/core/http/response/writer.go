package response

import (
	"encoding/json"
	"net/http"

	responsedto "grovehub/libs/core/dto/response"
	coreerrors "grovehub/libs/core/errors"
)

func WriteSuccess[T any](writer http.ResponseWriter, statusCode int, data T) error {
	return writeJSON(writer, statusCode, responsedto.NewSuccess(data))
}

func WriteError(writer http.ResponseWriter, statusCode int, definition coreerrors.Definition) error {
	return WriteErrorWithMetadata(writer, statusCode, definition, struct{}{})
}

func WriteErrorWithMetadata[M any](writer http.ResponseWriter, statusCode int, definition coreerrors.Definition, metadata M) error {
	return writeJSON(
		writer,
		statusCode,
		responsedto.NewError(definition.Code(), definition.Message(), statusCode, metadata),
	)
}

func writeJSON[T any](writer http.ResponseWriter, statusCode int, payload T) error {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(statusCode)
	return json.NewEncoder(writer).Encode(payload)
}
