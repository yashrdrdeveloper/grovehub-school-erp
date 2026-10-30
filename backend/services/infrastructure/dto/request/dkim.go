package request

type EnsureDKIM struct {
	Domain   string `json:"domain"`
	Selector string `json:"selector"`
}
