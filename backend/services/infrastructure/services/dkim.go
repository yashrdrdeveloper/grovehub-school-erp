package services

import "context"

type EnsureDKIMInput struct {
	Domain   string
	Selector string
}

type DKIMIdentity struct {
	Domain    string
	Selector  string
	PublicKey string
}

type DKIMService interface {
	Ensure(ctx context.Context, input EnsureDKIMInput) (DKIMIdentity, error)
}
