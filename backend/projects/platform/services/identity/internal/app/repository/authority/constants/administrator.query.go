package constants

// Each required permission has a code-owned feature binding. Counting distinct
// codes per account implements union across active roles, not across accounts.
const CountAdministrators = `SELECT count(*) FROM (
    SELECT a.account_id
    FROM identity.accounts a
    JOIN identity.account_roles ar ON ar.tenant_id = a.tenant_id AND ar.account_id = a.account_id
    JOIN identity.roles r ON r.tenant_id = ar.tenant_id AND r.role_id = ar.role_id
    JOIN identity.role_permissions rp ON rp.tenant_id = r.tenant_id AND rp.role_id = r.role_id
    JOIN identity.permissions p ON p.permission_id = rp.permission_id
    JOIN identity.features f ON f.feature_id = p.feature_id
    JOIN identity.tenant_features tf ON tf.tenant_id = a.tenant_id AND tf.feature_id = f.feature_id
    WHERE a.tenant_id = ? AND a.status = 'ACTIVE' AND r.status = 'ACTIVE'
        AND p.active AND f.active AND tf.enabled AND f.code = ?
        AND p.code IN ?
    GROUP BY a.account_id
    HAVING count(DISTINCT p.code) = ?
) AS administrators`
