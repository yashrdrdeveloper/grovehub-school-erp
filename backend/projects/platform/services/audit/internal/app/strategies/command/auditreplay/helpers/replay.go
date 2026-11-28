package helpers

import (
	"context"

	"gorm.io/gorm"

	activityrepo "grovehub/services/audit/internal/app/repository/activity/impl"
	"grovehub/services/audit/internal/app/services/replay"
	replayimpl "grovehub/services/audit/internal/app/services/replay/impl"
	"grovehub/services/audit/internal/app/strategies/command"
	"grovehub/services/audit/internal/app/strategies/command/auditreplay/constants"
	"grovehub/services/infrastructure"
)

type ReplayOperation func(context.Context, *gorm.DB, *gorm.DB, replay.Request) (replay.Result, error)

func Replay(ctx context.Context, auditDB, businessDB *gorm.DB, request replay.Request) (replay.Result, error) {
	return replayimpl.New(auditDB, infrastructure.NewReplaySource(businessDB), activityrepo.New(auditDB)).Replay(ctx, request)
}

func Execute(ctx context.Context, database command.Database[*gorm.DB], operation ReplayOperation, options Options) (replay.Result, error) {
	db, closeDB, err := database.Business(ctx)
	if err != nil {
		return replay.Result{}, constants.ErrBusinessConnection
	}
	defer closeDB()

	current, err := database.CurrentName(ctx, db)
	if err != nil || current != options.Database {
		return replay.Result{}, constants.ErrBusinessDatabaseMismatch
	}

	auditDB, closeAudit, err := database.Audit(ctx, db)
	if err != nil {
		return replay.Result{}, err
	}
	defer closeAudit()

	current, err = database.CurrentName(ctx, auditDB)
	if err != nil || current != options.AuditDatabase {
		return replay.Result{}, constants.ErrAuditDatabaseMismatch
	}

	result, err := operation(ctx, auditDB, db, options.Request)
	if err != nil {
		return replay.Result{}, constants.ErrReplay
	}

	return result, nil
}
