package impl

import (
	"context"

	"gorm.io/gorm"

	"grovehub/services/audit/internal/app/entity"
	"grovehub/services/audit/internal/app/repository/operator"
	"grovehub/services/audit/internal/app/repository/operator/constants"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) operator.Repository {
	return &repository{db: db}
}

func (r *repository) InsertActivity(ctx context.Context, a entity.OperatorActivity) error {
	args := []any{a.EventID, a.ActorID, a.ResourceID, a.Operation, a.OccurredAt, a.CorrelationID}
	return r.insertImmutable(ctx, constants.InsertActivity, constants.IdenticalActivity, args)
}

func (r *repository) InsertReceipt(ctx context.Context, a entity.SchoolMutationReceipt) error {
	args := []any{a.ExecutionID, a.SchoolID, a.PreviousVersion, a.CorrelationID, a.OccurredAt, a.Operation}
	return r.insertImmutable(ctx, constants.InsertReceipt, constants.IdenticalReceipt, args)
}

func (r *repository) insertImmutable(ctx context.Context, insert, identical string, args []any) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Exec(insert, args...)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 1 {
			return nil
		}

		var same bool
		if err := tx.Raw(identical, args...).Scan(&same).Error; err != nil {
			return err
		}
		if !same {
			return operator.ErrConflict
		}

		return nil
	})
}

func (r *repository) SchoolExists(ctx context.Context, school string) (bool, error) {
	var exists bool
	err := r.db.WithContext(ctx).Raw(constants.SchoolExists, school).Scan(&exists).Error
	return exists, err
}

func (r *repository) SearchActivities(ctx context.Context, limit, offset int) ([]entity.OperatorActivity, error) {
	records := make([]entity.OperatorActivity, 0)
	rows, err := r.db.WithContext(ctx).Raw(constants.SearchActivities, limit, offset).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a entity.OperatorActivity
		if err := rows.Scan(&a.EventID, &a.ActorID, &a.ResourceID, &a.Operation, &a.OccurredAt, &a.CorrelationID); err != nil {
			return nil, err
		}
		records = append(records, a)
	}

	return records, rows.Err()
}

func (r *repository) SearchReceipts(ctx context.Context, school, operation string, limit, offset int) ([]entity.SchoolMutationReceipt, error) {
	query := constants.SearchReceipts
	args := []any{}

	if school != "" {
		query += constants.SchoolFilter
		args = append(args, school)
	}
	if operation != "" {
		query += constants.OperationFilter
		args = append(args, operation)
	}

	query += constants.ReceiptOrder
	args = append(args, limit, offset)
	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]entity.SchoolMutationReceipt, 0)
	for rows.Next() {
		var a entity.SchoolMutationReceipt
		if err := rows.Scan(&a.ExecutionID, &a.SchoolID, &a.PreviousVersion, &a.CorrelationID, &a.OccurredAt, &a.Operation); err != nil {
			return nil, err
		}
		records = append(records, a)
	}

	return records, rows.Err()
}
