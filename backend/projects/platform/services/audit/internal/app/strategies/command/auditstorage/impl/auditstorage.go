package impl

import (
	"context"
	"encoding/json"
	"io"

	"gorm.io/gorm"

	"grovehub/services/audit/internal/app/services/storage"
	"grovehub/services/audit/internal/app/strategies/command"
	"grovehub/services/audit/internal/app/strategies/command/auditstorage/helpers"
	"grovehub/services/audit/internal/app/strategies/command/constants"
	shared "grovehub/services/audit/internal/app/strategies/command/helpers"
)

// @Strategy
type strategy struct {
	database command.Database[*gorm.DB]
	inspect  helpers.Inspector
}

func New(db command.Database[*gorm.DB]) command.Command {
	return NewWithInspector(db, storage.InspectStorage)
}

func NewWithInspector(db command.Database[*gorm.DB], operation helpers.Inspector) command.Command {
	return &strategy{database: db, inspect: operation}
}

func (s *strategy) Run(args []string, output, diagnostics io.Writer) int {
	return shared.Run(diagnostics, "", func() error {
		options, err := helpers.Parse(args, diagnostics)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), constants.Timeout)
		defer cancel()

		result, err := helpers.Execute(ctx, s.database, s.inspect, options)
		if err != nil {
			return err
		}

		return json.NewEncoder(output).Encode(result)
	})
}
