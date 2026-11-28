package controller
import "net/http"

// @Routes /
type HealthController interface {
 // @GET /health
 Health(http.ResponseWriter,*http.Request)
 // @GET /ready
 Readiness(http.ResponseWriter,*http.Request)
}
