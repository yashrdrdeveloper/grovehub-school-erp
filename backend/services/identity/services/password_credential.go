package services

import "context"

type PasswordCredentialService interface {
	Set(ctx context.Context, accountId string, password string) error
	Verify(ctx context.Context, accountId string, password string) (bool, error)
	Replace(ctx context.Context, accountId string, password string) error
}
