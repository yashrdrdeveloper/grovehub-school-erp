package impl

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"grovehub/services/mailProvider/entity"
	mailproviderrepository "grovehub/services/mailProvider/repository"
)

const uniqueSenderEmailConstraint = "uq_mail_provider_senders_email"

type SenderRepository struct {
	db *gorm.DB
}

func NewSenderRepository(db *gorm.DB) mailproviderrepository.SenderRepository { return &SenderRepository{db: db} }

func (repository *SenderRepository) Create(ctx context.Context, sender entity.Sender) error {
	return mapSenderWriteError(repository.db.WithContext(ctx).Create(&sender).Error)
}

func (repository *SenderRepository) FindAll(ctx context.Context) ([]entity.Sender, error) {
	senders := make([]entity.Sender, 0)
	if err := repository.db.WithContext(ctx).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).
		Find(&senders).Error; err != nil {
		return nil, err
	}
	return senders, nil
}

func (repository *SenderRepository) FindById(ctx context.Context, senderId entity.SenderId) (entity.Sender, error) {
	var sender entity.Sender
	if err := repository.db.WithContext(ctx).Where(&entity.Sender{SenderId: senderId}).Take(&sender).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Sender{}, mailproviderrepository.ErrSenderNotFound
		}
		return entity.Sender{}, err
	}
	return sender, nil
}

func (repository *SenderRepository) FindByEmail(ctx context.Context, email string) (entity.Sender, error) {
	var sender entity.Sender
	if err := repository.db.WithContext(ctx).Where(&entity.Sender{Email: email}).Take(&sender).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Sender{}, mailproviderrepository.ErrSenderNotFound
		}
		return entity.Sender{}, err
	}
	return sender, nil
}

func (repository *SenderRepository) Update(ctx context.Context, sender entity.Sender) error {
	expectedVersion := sender.GetVersion() - 1
	result := repository.db.WithContext(ctx).
		Model(&entity.Sender{}).
		Where(&entity.Sender{SenderId: sender.GetSenderId()}).
		Where("version = ?", expectedVersion).
		Select("Status", "Version", "UpdatedAt").
		Updates(&sender)
	if result.Error != nil {
		return mapSenderWriteError(result.Error)
	}
	if result.RowsAffected == 0 {
		return mailproviderrepository.ErrSenderVersionConflict
	}
	return nil
}

func mapSenderWriteError(err error) error {
	if err == nil {
		return nil
	}
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.Code == "23505" && postgresError.ConstraintName == uniqueSenderEmailConstraint {
		return mailproviderrepository.ErrSenderAlreadyExists
	}
	return err
}
