// Package replace provides an unwired guarded assignment persistence operation.
// The future approval executor must revalidate actor authority and payload in
// the same tenant transaction before calling this operation.
package replace

import (
	"context"

	"grovehub/services/identity/internal/app/repository/accountrole"
	"grovehub/services/identity/internal/app/repository/accountrole/helpers"
	"grovehub/services/identity/internal/app/repository/authority"
	"grovehub/services/identity/internal/app/services/authorization"
)

type Operation struct {
	guard authority.Guard
	writer accountrole.Writer
}

func New(guard authority.Guard, writer accountrole.Writer) *Operation {
	return &Operation{guard: guard, writer: writer}
}

func (operation *Operation) Execute(ctx context.Context, request accountrole.Replacement) (accountrole.Assignment, error) {
	request, err := helpers.NormalizeReplacement(request)
	if err != nil {
		return accountrole.Assignment{}, err
	}
	if operation.guard == nil || operation.writer == nil {
		return accountrole.Assignment{}, authorization.ErrUnavailable
	}
	var assignment accountrole.Assignment
	err = operation.guard.Run(ctx, request.TenantID, func(tx context.Context) error {
		var err error
		assignment, err = operation.writer.Replace(tx, request)
		return err
	})
	if err != nil {
		return accountrole.Assignment{}, err
	}
	return assignment, nil
}
