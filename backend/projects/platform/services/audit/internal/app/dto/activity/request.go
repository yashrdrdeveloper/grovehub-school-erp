package activity

import coresearch "grovehub/libs/core/search"

type Search struct {
	coresearch.Request
	Scope string `json:"scope,omitempty"`
}
