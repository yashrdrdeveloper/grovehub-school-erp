package person

// PersonByGID identifies the owning Person from the route path.
type PersonByGID struct {
	GID string `json:"gid" path:"gid" normalize:"trim,uuidv7" validate:"uuidv7"`
}

