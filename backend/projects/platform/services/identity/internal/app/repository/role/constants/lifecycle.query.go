package constants

const TransitionStatus = `UPDATE identity.roles SET status = ?, version = version + 1,
updated_by = ?, updated_at = now()
WHERE tenant_id = ? AND role_id = ? AND version = ? AND status = ?`

const EnabledFeature = `SELECT count(*) FROM identity.features f
JOIN identity.tenant_features tf ON tf.feature_id = f.feature_id
WHERE tf.tenant_id = ? AND f.code = ? AND f.active AND tf.enabled`
