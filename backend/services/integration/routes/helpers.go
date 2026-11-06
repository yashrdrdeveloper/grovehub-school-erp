package routes

import (
 "net/http"
 httptracing "grovehub/libs/core/http/tracing"
)

// handle registers one traced handler; route definitions live with their domain.
func handle(mux *http.ServeMux,method,path string,handler http.HandlerFunc){
 pattern:=method+" "+path
 mux.Handle(pattern,httptracing.ServiceHandler("integration",pattern,handler))
}
