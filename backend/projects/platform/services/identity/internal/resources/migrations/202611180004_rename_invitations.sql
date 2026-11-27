-- +goose Up
-- Preserve existing invitation delivery jobs while aligning the entity/table name.
ALTER TABLE identity.invitation_mail RENAME TO invitations;
ALTER INDEX IF EXISTS identity.idx_identity_invitation_mail_pending RENAME TO idx_identity_invitations_pending;

-- +goose Down
ALTER TABLE identity.invitations RENAME TO invitation_mail;
ALTER INDEX IF EXISTS identity.idx_identity_invitations_pending RENAME TO idx_identity_invitation_mail_pending;
