package events

const (
	ActivityOccurred = "audit.activity.occurred.v1"
	ActivityOccurredTopic = "audit_activity_occurred"
	AccountChanged = "identity.account.changed.v1"
	AccountChangedTopic = "identity_account_changed"
	AccountRegistrationChanged = "identity.account.registration.changed.v1"
	AccountRegistrationChangedTopic = "identity_account_registration_changed"
	AccountMetadataChanged = "identity.account.metadata.changed.v1"
	AccountMetadataChangedTopic = "identity_account_metadata_changed"
	AccountAuthorityChanged = "identity.account.authority.changed.v1"
	AccountAuthorityChangedTopic = "identity_account_authority_changed"
	RoleChanged = "identity.role.changed.v1"
	RoleChangedTopic = "identity_role_changed"
	RoleMetadataChanged = "identity.role.metadata.changed.v1"
	RoleMetadataChangedTopic = "identity_role_metadata_changed"
	ApprovalChanged = "identity.approval.changed.v1"
	ApprovalChangedTopic = "identity_approval_changed"
	InvitationDelivery = "identity.invitation.delivery.changed.v1"
	InvitationDeliveryTopic = "identity_invitation_delivery_changed"
)

// Topics returns the explicit subscription allowlist.
func Topics() map[string]string {
	return map[string]string{
		ActivityOccurredTopic: ActivityOccurred,
		AccountChangedTopic: AccountChanged,
		AccountRegistrationChangedTopic: AccountRegistrationChanged,
		AccountMetadataChangedTopic: AccountMetadataChanged,
		AccountAuthorityChangedTopic: AccountAuthorityChanged,
		RoleChangedTopic: RoleChanged,
		RoleMetadataChangedTopic: RoleMetadataChanged,
		ApprovalChangedTopic: ApprovalChanged,
		InvitationDeliveryTopic: InvitationDelivery,
	}
}

func Topic(eventType string) (string, bool) {
	for topic, kind := range Topics() {
		if kind == eventType {
			return topic, true
		}
	}
	return "", false
}
