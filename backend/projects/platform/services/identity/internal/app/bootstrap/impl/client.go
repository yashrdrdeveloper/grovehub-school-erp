package impl

import (
	"net/http"

	accountclient "grovehub/services/identity/client/account"
	clientimpl "grovehub/services/identity/internal/app/client/account/impl"
)

func (b *builder) newClient(controllers controllerGraph) accountclient.Client {
	internalMux := http.NewServeMux()
	controllers.RegisterHTTP(internalMux)
	return clientimpl.NewAccountClient(internalMux)
}
