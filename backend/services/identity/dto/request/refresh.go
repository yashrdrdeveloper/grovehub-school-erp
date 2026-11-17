package request

import corerequest "grovehub/libs/core/request"

type Refresh struct { RefreshToken string `json:"refreshToken"` }
type RefreshRequest struct {
    corerequest.Request
    RequestBody Refresh `validate:"valid"`
}
