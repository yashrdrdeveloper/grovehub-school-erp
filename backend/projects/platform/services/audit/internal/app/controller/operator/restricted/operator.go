package restricted

import (
	"net/http"

	"grovehub/services/audit/internal/app/controller/operator"
	"grovehub/services/audit/internal/app/controller/operator/constants"
)

func Register(mux *http.ServeMux, controller operator.Controller) {
	mux.HandleFunc(constants.Pattern, controller.Write)
	mux.HandleFunc(constants.SearchPattern, controller.Search)
	mux.HandleFunc(constants.SchoolCheckPattern, controller.SchoolCheck)
	mux.HandleFunc(constants.ReceiptPattern, controller.Receipt)
	mux.HandleFunc(constants.StatusReceiptPattern, controller.StatusReceipt)
	mux.HandleFunc(constants.ReceiptSearchPattern, controller.ReceiptSearch)
}
