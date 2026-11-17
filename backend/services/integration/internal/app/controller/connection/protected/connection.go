package protected

import "net/http"

// @Routes /connections
type ConnectionController interface {
	// @POST /
	Create(http.ResponseWriter, *http.Request)
	// @GET /{id}
	Get(http.ResponseWriter, *http.Request)
	// @GET /
	List(http.ResponseWriter, *http.Request)
	// @DELETE /{id}
	Disconnect(http.ResponseWriter, *http.Request)
	// @POST /{id}/oauth/start
	BeginOAuth(http.ResponseWriter, *http.Request)
	// @POST /{id}/test
	RunTest(http.ResponseWriter, *http.Request)
	// @GET /{id}/tests
	ListTests(http.ResponseWriter, *http.Request)
	// @POST /{id}/watch
	Watch(http.ResponseWriter, *http.Request)
	// @POST /{id}/actions/{action}
	ExecuteAction(http.ResponseWriter, *http.Request)
}
