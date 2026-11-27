package constants

const Module = "IDENTITY"

const InsertFeature = `INSERT INTO identity.features (feature_id, code, module, name, active)
VALUES (?, ?, ?, ?, false) ON CONFLICT (code) DO NOTHING`

const FindFeature = `SELECT feature_id::text AS id, module
FROM identity.features WHERE code = ? FOR SHARE`

const InsertPermission = `INSERT INTO identity.permissions (permission_id, code, module, name, feature_id, active)
VALUES (?, ?, ?, ?, ?, false) ON CONFLICT (code) DO NOTHING`

const FindPermission = `SELECT permission_id::text AS id, module, feature_id::text AS feature_id
FROM identity.permissions WHERE code = ? FOR SHARE`
