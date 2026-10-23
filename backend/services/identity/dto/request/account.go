package request

import foundationrequest "grovehub/services/foundation/dto/request"

type LoginIdentifier struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type BootstrapAccount struct {
	Person     foundationrequest.CreatePerson `json:"person"`
	Identifier LoginIdentifier                `json:"identifier"`
	Password   string                         `json:"password"`
}
