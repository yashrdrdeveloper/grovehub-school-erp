package impl

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"gorm.io/gorm"
	"grovehub/services/audit/internal/app/repository/rejection"
)

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) rejection.Repository {
	return &Repository{db: db}
}

func (repository *Repository) Record(ctx context.Context, eventType string, payload []byte, reason string) error {
	hash := sha256.New()
	hash.Write([]byte(eventType))
	hash.Write([]byte{0})
	hash.Write([]byte(reason))
	hash.Write([]byte{0})
	hash.Write(payload)
	return repository.db.WithContext(ctx).Exec("INSERT INTO audit.rejected_events(fingerprint,event_type,reason,payload_bytes) VALUES (?,?,?,?) ON CONFLICT (fingerprint) DO NOTHING", hex.EncodeToString(hash.Sum(nil)), eventType, reason, len(payload)).Error
}
