package impl

import (
 "context"
 "errors"
 "gorm.io/gorm"
)

type TenantSessions struct { db *gorm.DB }

func NewTenantSessions(db *gorm.DB) *TenantSessions { return &TenantSessions{db: db} }

// ListTenants uses a keyset cursor over tenants with potentially usable
// sessions. The worker wraps at the end so failed tenants are retried.
func (r *TenantSessions) ListTenants(ctx context.Context, after string, limit int) ([]string, error) {
 if limit < 1 || limit > 100 { return nil, errors.New("invalid tenant batch size") }
 var tenants []string
 err := r.db.WithContext(ctx).Raw(`
  SELECT DISTINCT a.tenant_id
  FROM identity.accounts a
  JOIN identity.sessions s ON s.account_id = a.account_id
  WHERE a.tenant_id > ? AND s.revoked_at IS NULL
    AND s.refresh_expires_at > now()
  ORDER BY a.tenant_id LIMIT ?`, after, limit).Scan(&tenants).Error
 return tenants, err
}

// RevokeBefore never uses an inequality that could revoke sessions issued
// after a delayed registry response. A repeated call is a no-op.
func (r *TenantSessions) RevokeBefore(ctx context.Context, tenant string, generation int64) (int64, error) {
 if tenant == "" || generation < 1 { return 0, errors.New("invalid tenant revocation cutoff") }
 result := r.db.WithContext(ctx).Exec(`
  UPDATE identity.sessions s SET revoked_at = now(), updated_at = now()
  FROM identity.accounts a
  WHERE s.account_id = a.account_id AND a.tenant_id = ?
    AND s.revoked_at IS NULL AND s.session_generation < ?`, tenant, generation)
 return result.RowsAffected, result.Error
}
