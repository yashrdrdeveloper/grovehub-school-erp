package external

import "net/http"

type Operation interface { http.Handler }
