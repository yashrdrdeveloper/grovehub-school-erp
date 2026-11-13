package impl

import (
	"net/http"

	foundationclient "grovehub/services/foundation/client"
)

// personClient dispatches annotated requests through the owning controller's
// HTTP route pipeline in-process. No TCP connection or service proxy is used.
type personClient struct {
	handler http.Handler
}

func NewPersonClient(handler http.Handler) foundationclient.Client {
	return &personClient{handler: handler}
}
