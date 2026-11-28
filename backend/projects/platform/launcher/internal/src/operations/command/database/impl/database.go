package impl

import (
	"context"

	"gorm.io/gorm"

	bootstrapimpl "grovehub/launcher/internal/src/bootstrap/impl"
	"grovehub/libs/platform/command"
)

type BusinessOpener func(context.Context) (*gorm.DB, error)
type AuditOpener func(context.Context, *gorm.DB) (*gorm.DB, func(), error)

// @Service
type postgres struct {
	openBusiness BusinessOpener
	openAudit    AuditOpener
}

func New() command.Database[*gorm.DB] {
	return NewWithOpeners(bootstrapimpl.NewPostgresDB, bootstrapimpl.NewAuditPostgresDB)
}

func NewWithOpeners(business BusinessOpener, audit AuditOpener) command.Database[*gorm.DB] {
	return &postgres{openBusiness: business, openAudit: audit}
}

func (p *postgres) Business(ctx context.Context) (*gorm.DB, func(), error) {
	db, err := p.openBusiness(ctx)
	if err != nil {
		return nil, nil, err
	}

	raw, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	return db, func() { _ = raw.Close() }, nil
}

func (p *postgres) Audit(ctx context.Context, shared *gorm.DB) (*gorm.DB, func(), error) {
	return p.openAudit(ctx, shared)
}

func (p *postgres) CurrentName(ctx context.Context, db *gorm.DB) (string, error) {
	var current string
	err := db.WithContext(ctx).Raw("SELECT current_database()").Scan(&current).Error
	return current, err
}
