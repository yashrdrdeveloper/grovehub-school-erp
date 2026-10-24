package response

type Registration struct {
	AccountId string `json:"accountId"`
	PersonGid string `json:"personGid"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
	Status    string `json:"status"`
}
