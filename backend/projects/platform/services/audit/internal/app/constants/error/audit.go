package errorcode

import (
	"net/http"
	coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
)

var InvalidRequest = coreerrors.NewDefinition("AUDIT_REQUEST_INVALID","Invalid audit search",http.StatusBadRequest)
var Forbidden = coreerrors.NewDefinition("AUDIT_FORBIDDEN","Access denied",http.StatusForbidden)
var Internal = coreerrors.NewDefinition("AUDIT_INTERNAL","Audit operation failed",http.StatusInternalServerError)
