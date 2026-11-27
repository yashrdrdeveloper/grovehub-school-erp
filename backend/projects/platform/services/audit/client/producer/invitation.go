package producer

import "time"

const InvitationDelivery = "identity.invitation.delivery.changed.v1"

// InvitationEvent contains only delivery lifecycle metadata, never mail content.
type InvitationEvent struct {
	TenantID string `json:"tenantId"`
	AccountID string `json:"accountId"`
	Purpose string `json:"purpose"`
	Operation string `json:"operation"`
	Attempt int `json:"attempt"`
	Reason string `json:"reason,omitempty"`
	NextAttemptAt *time.Time `json:"nextAttemptAt,omitempty"`
}
