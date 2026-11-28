package impl

import (
	"context"
	"encoding/json"
	"io"

	"gorm.io/gorm"

	bootstrap "grovehub/services/audit/internal/app/bootstrap/constants"
	"grovehub/services/audit/internal/app/strategies/command"
	"grovehub/services/audit/internal/app/strategies/command/auditreplay/helpers"
	"grovehub/services/audit/internal/app/strategies/command/constants"
	shared "grovehub/services/audit/internal/app/strategies/command/helpers"
)

// @Strategy
type strategy struct {
	database command.Database[*gorm.DB]
	replay helpers.ReplayOperation
}

func New(db command.Database[*gorm.DB]) command.Command {
	return NewWithReplay(db, helpers.Replay)
}

func NewWithReplay(db command.Database[*gorm.DB], operation helpers.ReplayOperation) command.Command {
	return &strategy{database: db, replay: operation}
}

func (s *strategy) Run(args []string, output, diagnostics io.Writer) int {
	return shared.Run(diagnostics, bootstrap.AuditReplay + ":", func() error {
		options, err := helpers.Parse(args, diagnostics)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), constants.Timeout)
		defer cancel()

		result, err := helpers.Execute(ctx, s.database, s.replay, options)
		if err != nil {
			return err
		}

		return json.NewEncoder(output).Encode(result)
	})
}
