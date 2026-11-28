// Package entity provides a generic, reflection-free builder for passive entities.
package entity

// Builder creates a value of T using typed mutations, without generated methods.
type Builder[T any] struct {
	value T
}

// New returns a builder initialized with the zero value of T.
func New[T any]() Builder[T] {
	return Builder[T]{}
}

// From starts a builder with an existing entity value.
func From[T any](value T) Builder[T] {
	return Builder[T]{value: value}
}

// With applies a typed update to the builder's own copy and returns the new builder.
// This does not perform business validation, persistence, or lifecycle transitions.
func (builder Builder[T]) With(change func(*T)) Builder[T] {
	if change != nil {
		change(&builder.value)
	}
	return builder
}

// Build returns the current entity value.
func (builder Builder[T]) Build() T {
	return builder.value
}
