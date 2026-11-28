package constants

const Available = `SELECT p.code, COALESCE(f.code, '') AS feature
FROM identity.permissions p
LEFT JOIN identity.features f ON f.feature_id = p.feature_id
LEFT JOIN identity.tenant_features tf ON tf.feature_id = f.feature_id AND tf.tenant_id = ?
WHERE p.module = 'IDENTITY' AND p.active
AND (p.feature_id IS NULL OR (f.module = 'IDENTITY' AND f.active AND tf.enabled))
ORDER BY p.code`
