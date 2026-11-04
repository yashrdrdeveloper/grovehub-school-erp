package response

type Register struct {
	AccountId string `json:"accountId"`
	PersonGid string `json:"personGid"`
	GID       string `json:"gid"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
	Status    string `json:"status"`
}
