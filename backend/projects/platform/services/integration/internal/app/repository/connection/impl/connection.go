package impl

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	auditclient "grovehub/services/audit/client/producer"
	"time"

	coredatabase "grovehub/libs/core/database"
	"grovehub/services/integration/internal/app/entity"
	integrationrepository "grovehub/services/integration/internal/app/repository/connection"
)

type ConnectionRepository struct {
	db    *gorm.DB
	audit auditclient.Client
}

func NewConnectionRepository(db *gorm.DB, clients ...auditclient.Client) integrationrepository.ConnectionRepository {
	repository := &ConnectionRepository{db: db}
	if len(clients) > 0 {
		repository.audit = clients[0]
	}
	return repository
}

func (r *ConnectionRepository) Create(ctx context.Context, c entity.Connection) error {
	return coredatabase.RunTransaction(ctx, r.db, func(tx context.Context) error {
		if err := coredatabase.TransactionDB(tx, r.db).Create(&c).Error; err != nil {
			return err
		}
		return r.recordActivity(tx, c, "CREATED", c.CreatedAt)
	})
}

func (r *ConnectionRepository) Get(ctx context.Context, tenantID, id string) (entity.Connection, error) {
	var c entity.Connection
	err := coredatabase.TransactionDB(ctx, r.db).Where("tenant_id = ? AND id = ?", tenantID, id).Take(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c, integrationrepository.ErrConnectionNotFound
	}
	return c, err
}

func (r *ConnectionRepository) GetByID(ctx context.Context, id string) (entity.Connection, error) {
	var c entity.Connection
	err := coredatabase.TransactionDB(ctx, r.db).Where("id = ?", id).Take(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c, integrationrepository.ErrConnectionNotFound
	}
	return c, err
}

func (r *ConnectionRepository) List(ctx context.Context, tenantID, ownerGID string) ([]entity.Connection, error) {
	rows := make([]entity.Connection, 0)
	err := coredatabase.TransactionDB(ctx, r.db).Where("tenant_id = ? AND owner_gid = ?", tenantID, ownerGID).Order("created_at DESC, id DESC").Limit(200).Find(&rows).Error
	return rows, err
}

func (r *ConnectionRepository) Disconnect(ctx context.Context, tenantID string, id string) (entity.Connection, error) {
	var entry entity.Connection
	err := coredatabase.RunTransaction(ctx, r.db, func(tx context.Context) error {
		db := coredatabase.TransactionDB(tx, r.db)
		if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("tenant_id = ? AND id = ?", tenantID, id).Take(&entry).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return integrationrepository.ErrConnectionNotFound
			}
			return err
		}
		if entry.Status == entity.StatusDisconnected {
			return nil
		}
		now := time.Now().UTC()
		if err := db.Model(&entity.Connection{}).Where("tenant_id = ? AND id = ?", tenantID, id).Updates(map[string]any{"status": entity.StatusDisconnected, "credential_ref": nil, "updated_at": now}).Error; err != nil {
			return err
		}
		if err := db.Where("connection_id = ?", id).Delete(&entity.OAuthCredential{}).Error; err != nil {
			return err
		}
		if err := db.Where("connection_id = ?", id).Delete(&entity.OAuthState{}).Error; err != nil {
			return err
		}
		entry.Status = entity.StatusDisconnected
		entry.CredentialRef = nil
		entry.UpdatedAt = now
		return r.recordActivity(tx, entry, "DISCONNECTED", now)
	})
	if err != nil {
		return entity.Connection{}, err
	}
	return entry, nil
}

func (r *ConnectionRepository) SaveOAuthState(ctx context.Context, state entity.OAuthState) error {
	return r.db.WithContext(ctx).Create(&state).Error
}

func (r *ConnectionRepository) ConsumeOAuthState(ctx context.Context, hash string, now time.Time) (entity.OAuthState, error) {
	var state entity.OAuthState
	err := r.db.WithContext(ctx).Clauses(clause.Returning{}).
		Where("hash = ? AND expires_at > ?", hash, now).
		Delete(&state).Error
	if err != nil {
		return entity.OAuthState{}, err
	}
	if state.Hash == "" {
		return entity.OAuthState{}, integrationrepository.ErrOAuthStateInvalid
	}
	return state, nil
}

func (r *ConnectionRepository) Activate(ctx context.Context, connectionID string, provider string, externalID string, ciphertext []byte) error {
	return coredatabase.RunTransaction(ctx, r.db, func(tx context.Context) error {
		db := coredatabase.TransactionDB(tx, r.db)
		now := time.Now().UTC()
		result := db.Model(&entity.Connection{}).Where("id = ? AND provider = ? AND status = ?", connectionID, provider, entity.StatusPendingSetup).Updates(map[string]any{"external_id": externalID, "status": entity.StatusActive, "updated_at": now})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return integrationrepository.ErrConnectionNotFound
		}
		if err := db.Create(&entity.OAuthCredential{ConnectionID: connectionID, Ciphertext: ciphertext, UpdatedAt: now}).Error; err != nil {
			return err
		}
		entry, err := r.GetByID(tx, connectionID)
		if err != nil {
			return err
		}
		return r.recordActivity(tx, entry, "ACTIVATED", now)
	})
}

func (r *ConnectionRepository) Credential(ctx context.Context, connectionID string) ([]byte, error) {
	var credential entity.OAuthCredential
	err := r.db.WithContext(ctx).Where("connection_id = ?", connectionID).Take(&credential).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, integrationrepository.ErrConnectionNotFound
	}
	return credential.Ciphertext, err
}

func (r *ConnectionRepository) UpdateCredential(ctx context.Context, connectionID string, ciphertext []byte) error {
	return r.db.WithContext(ctx).Model(&entity.OAuthCredential{}).Where("connection_id = ?", connectionID).
		Updates(map[string]any{"ciphertext": ciphertext, "updated_at": time.Now().UTC()}).Error
}

func (r *ConnectionRepository) SaveTest(ctx context.Context, test entity.ConnectionTest) error {
	return coredatabase.RunTransaction(ctx, r.db, func(tx context.Context) error {
		db := coredatabase.TransactionDB(tx, r.db)
		if err := db.Create(&test).Error; err != nil {
			return err
		}
		if err := db.Model(&entity.Connection{}).Where("id = ? AND (last_tested_at IS NULL OR last_tested_at <= ?)", test.ConnectionID, test.CompletedAt).Updates(map[string]any{"last_tested_at": test.CompletedAt, "last_test_status": test.Status, "updated_at": time.Now().UTC()}).Error; err != nil {
			return err
		}
		entry, err := r.GetByID(tx, test.ConnectionID)
		if err != nil {
			return err
		}
		return r.recordActivity(tx, entry, "TESTED", test.CompletedAt)
	})
}

func (r *ConnectionRepository) ListTests(ctx context.Context, connectionID string) ([]entity.ConnectionTest, error) {
	tests := make([]entity.ConnectionTest, 0)
	err := r.db.WithContext(ctx).Where("connection_id = ?", connectionID).
		Order("completed_at DESC, id DESC").Limit(100).Find(&tests).Error
	return tests, err
}

func (r *ConnectionRepository) SystemMailSender(ctx context.Context, tenantID string) (entity.Connection, error) {
	var sender entity.Connection
	err := coredatabase.TransactionDB(ctx, r.db).
		Table("integration.connections AS c").
		Joins("JOIN integration.system_mail_senders AS s ON s.connection_id = c.id AND s.tenant_id = c.tenant_id").
		Where("s.tenant_id = ? AND c.status = ?", tenantID, entity.StatusActive).
		Take(&sender).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Connection{}, integrationrepository.ErrConnectionNotFound
	}
	return sender, err
}

func (r *ConnectionRepository) ConfigureSystemMailSender(ctx context.Context, tenantID string, connectionID string) error {
	return coredatabase.RunTransaction(ctx, r.db, func(tx context.Context) error {
		result := coredatabase.TransactionDB(tx, r.db).Exec(
			`INSERT INTO integration.system_mail_senders (tenant_id, connection_id, updated_at)
		SELECT c.tenant_id, c.id, CURRENT_TIMESTAMP
		FROM integration.connections AS c
		WHERE c.tenant_id = ? AND c.id = ? AND c.status = ?
		ON CONFLICT (tenant_id) DO UPDATE
		SET connection_id = EXCLUDED.connection_id, updated_at = EXCLUDED.updated_at`,
			tenantID, connectionID, entity.StatusActive,
		)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return integrationrepository.ErrConnectionNotFound
		}
		entry, err := r.Get(tx, tenantID, connectionID)
		if err != nil {
			return err
		}
		return r.recordActivity(tx, entry, "SYSTEM_MAIL_SENDER_CONFIGURED", time.Now().UTC())
	})
}
