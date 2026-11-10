package controller

import "net/http"

// @Routes /integrations
type ConnectionController interface {
    // @POST /connections
    Create(http.ResponseWriter, *http.Request)
    // @GET /connections/{id}
    Get(http.ResponseWriter, *http.Request)
    // @GET /connections
    List(http.ResponseWriter, *http.Request)
    // @DELETE /connections/{id}
    Disconnect(http.ResponseWriter, *http.Request)
    // @POST /connections/{id}/oauth/start
    BeginOAuth(http.ResponseWriter, *http.Request)
    // @GET /connections/oauth/{provider}/callback
    FinishOAuth(http.ResponseWriter, *http.Request)
    // @POST /connections/{id}/test
    RunTest(http.ResponseWriter, *http.Request)
    // @GET /connections/{id}/tests
    ListTests(http.ResponseWriter, *http.Request)
    // @POST /connections/{id}/watch
    Watch(http.ResponseWriter, *http.Request)
    // @POST /webhooks/gmail
    GmailPush(http.ResponseWriter, *http.Request)
    // @POST /connections/{id}/actions/{action}
    ExecuteAction(http.ResponseWriter, *http.Request)
}
