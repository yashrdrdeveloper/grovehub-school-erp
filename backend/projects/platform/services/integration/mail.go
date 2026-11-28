package integration

import (
	"net/http"
	"gorm.io/gorm"
	"grovehub/services/identity/client/invitation"
	"grovehub/services/integration/client/mail"
	clientimpl "grovehub/services/integration/internal/app/client/mail/impl"
	controller "grovehub/services/integration/internal/app/controller/mail/restricted"
	service "grovehub/services/integration/internal/app/services/mail/impl"
)

// NewMailReconciliationClient composes Integration's read-only diagnostic
// with the Identity-owned invitation projection, without starting workers.
func NewMailReconciliationClient(db *gorm.DB, jobs invitation.DeliveryReader) mail.Client {
	return clientimpl.NewMailClient(service.New(db, jobs))
}

func RegisterMailReconciliationHTTP(mux *http.ServeMux, client mail.Client) {
	mux.HandleFunc("POST "+mail.ReconciliationPath,controller.SearchHandler(client))
}
