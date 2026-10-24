package request

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func DecodeJSON(request *http.Request, destination any) error {
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON value")
	}
	return nil
}
