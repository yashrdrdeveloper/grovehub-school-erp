package events

// Explicit allowlist: never subscribe to arbitrary opaque provider payloads.
func Topics() map[string]string {
	return map[string]string{
		"identity_account_changed": "identity.account.changed.v1",
		"identity_account_registration_changed": "identity.account.registration.changed.v1",
		"identity_account_metadata_changed": "identity.account.metadata.changed.v1",
		"identity_account_authority_changed": "identity.account.authority.changed.v1",
		"identity_role_changed": "identity.role.changed.v1",
		"identity_role_metadata_changed": "identity.role.metadata.changed.v1",
		"identity_approval_changed": "identity.approval.changed.v1",
		"identity_invitation_delivery_changed": "identity.invitation.delivery.changed.v1",
	}
}

func Topic(eventType string) (string, bool) {
	for topic, kind := range Topics() {
		if kind == eventType { return topic, true }
	}
	return "", false
}
