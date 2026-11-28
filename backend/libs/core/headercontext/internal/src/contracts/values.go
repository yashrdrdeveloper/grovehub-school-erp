package contracts

// Values is a small immutable copy of trusted header metadata.
type Values struct {
	IdempotencyKey string
	CorrelationID  string
	RequestID      string
}
