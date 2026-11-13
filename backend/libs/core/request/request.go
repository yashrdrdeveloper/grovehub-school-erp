// Package request defines the transport-independent typed request consumed by
// controllers, domain services and generated in-process clients.
package request

// Empty marks a path, query or body source absent from a request contract.
type Empty struct{}

// Request carries validated path parameters, query parameters and a body.
// It is an in-process contract, not an HTTP or JSON wire envelope.
type Request[Path, Query, Body any] struct {
	Path Path
	Query Query
	Body Body
}
