package request

import foundationrequest "grovehub/services/foundation/dto/request"

type Register struct {
	foundationrequest.CreatePerson
	Password string `json:"password"`
}
