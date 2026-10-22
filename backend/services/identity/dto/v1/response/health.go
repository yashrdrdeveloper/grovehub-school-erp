package response

type HealthData struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}
