package mail

import (
	"grovehub/services/integration/client/mail"
)

// Service is the Integration-owned reconciliation contract.
type Service interface { mail.Client }
