package response

type DKIMIdentity struct {
	Domain    string `json:"domain"`
	Selector  string `json:"selector"`
	PublicKey string `json:"publicKey"`
}
