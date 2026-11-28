-- +goose Up
-- Provider-owned capabilities only. No workload or permission grant is seeded.
INSERT INTO registry.capabilities(code, scope) VALUES
    ('registry.project.read', 'PROJECT'),
    ('registry.project.create', 'PLATFORM'),
    ('registry.project.status', 'PROJECT'),
    ('registry.tenant.read', 'TENANT'),
    ('registry.tenant.create', 'PROJECT'),
    ('registry.tenant.status', 'TENANT'),
    ('registry.workload.read', 'PLATFORM'),
    ('registry.workload.create', 'PLATFORM'),
    ('registry.workload.status', 'PLATFORM'),
    ('registry.key.register', 'PLATFORM'),
    ('registry.key.revoke', 'PLATFORM'),
    ('registry.grant.create', 'PLATFORM'),
    ('registry.grant.revoke', 'PLATFORM');

-- +goose Down
DELETE FROM registry.capabilities WHERE code IN (
    'registry.project.read', 'registry.project.create', 'registry.project.status',
    'registry.tenant.read', 'registry.tenant.create', 'registry.tenant.status',
    'registry.workload.read', 'registry.workload.create', 'registry.workload.status',
    'registry.key.register', 'registry.key.revoke',
    'registry.grant.create', 'registry.grant.revoke'
);
