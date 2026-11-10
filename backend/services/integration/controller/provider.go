package controller

import "net/http"

// @Routes /integrations
type ProviderController interface {
    // @GET /providers
    Providers(http.ResponseWriter, *http.Request)
    // @GET /providers/{provider}
    Provider(http.ResponseWriter, *http.Request)
}
