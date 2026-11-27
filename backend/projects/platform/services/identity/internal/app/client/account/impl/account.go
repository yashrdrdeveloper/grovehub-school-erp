package impl

import (
	"net/http"

	accountclient "grovehub/services/identity/client/account"
)

// identityAccountClient uses the generated RequestLine dispatch through the
// controller and REST validation pipeline, matching Foundation's client.
type identityAccountClient struct {
	handler http.Handler
}

func NewAccountClient(handler http.Handler) accountclient.Client {
	return &identityAccountClient{handler: handler}
}
