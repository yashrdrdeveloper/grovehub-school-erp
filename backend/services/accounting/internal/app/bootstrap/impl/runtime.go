package impl

import (
 "net/http"
 "grovehub/services/accounting/internal/app/bootstrap"
)

type runtime struct{}

var _ bootstrap.Runtime = (*runtime)(nil)

// RegisterHTTP intentionally exposes no routes until authorization is configured.
func (module *runtime) RegisterHTTP(mux *http.ServeMux) {}

func (module *runtime) Close() {}
