package events

const (
	ActivityOccurred = "audit.activity.occurred.v1"
	InvitationDelivery = "identity.invitation.delivery.changed.v1"
)

// Explicit allowlist: never subscribe to arbitrary opaque provider payloads.
func Topics() map[string]string {
	return map[string]string{
		"audit_activity_occurred": ActivityOccurred,
		"identity_account_changed": "identity.account.changed.v1",
		"identity_account_registration_changed": "identity.account.registration.changed.v1",
		"identity_account_metadata_changed": "identity.account.metadata.changed.v1",
		"identity_account_authority_changed": "identity.account.authority.changed.v1",
		"identity_role_changed": "identity.role.changed.v1",
		"identity_role_metadata_changed": "identity.role.metadata.changed.v1",
		"identity_approval_changed": "identity.approval.changed.v1",
		"identity_invitation_delivery_changed": InvitationDelivery,
	}
}

func Topic(eventType string) (string, bool) {
	for topic, kind := range Topics() {
		if kind == eventType { return topic, true }
	}
	return "", false
}
