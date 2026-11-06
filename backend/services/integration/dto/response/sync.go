package response

type SyncResult struct {
 ConnectionID string `json:"connectionId"`
 Mode string `json:"mode"`
 Processed int `json:"processed"`
 Deleted int `json:"deleted"`
 Cursor string `json:"cursor"`
}
