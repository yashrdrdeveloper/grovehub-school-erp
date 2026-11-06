package response
import "time"
type Watch struct{
 ConnectionID string `json:"connectionId"`
 ExpiresAt time.Time `json:"expiresAt"`
}
