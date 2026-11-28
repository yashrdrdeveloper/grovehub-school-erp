package health

type Response struct {
	Status string `json:"status"`
}

type Service interface {
	Health() Response
	Readiness() (Response, bool)
}
