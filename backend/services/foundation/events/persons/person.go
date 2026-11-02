package persons

const (
	CreatedTopic = "foundation_person_created"
	CreatedType  = "foundation.person.created.v1"
	Source       = "foundation"
)

type Created struct {
	GID     string `json:"gid"`
	Status  string `json:"status"`
	Version int64  `json:"version"`
}
