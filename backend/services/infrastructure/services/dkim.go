package services

import (
	"context"
	"errors"
)

var (
	ErrInvalidDKIMIdentity  = errors.New("invalid DKIM identity")
	ErrDKIMSelectorConflict = errors.New("DKIM selector does not match the persisted identity")
)

type DKIMIdentity struct {
	Domain    string `json:"domain"`
	Selector  string `json:"selector"`
	PublicKey string `json:"publicKey"`
}

type DKIMService interface {
	Ensure(ctx context.Context, domain string, selector string) (DKIMIdentity, error)
}
