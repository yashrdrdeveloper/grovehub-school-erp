package restricted

import (
	"context"
	"errors"
	"net/http"
	"time"

	coregateway "grovehub/libs/core/gateway"
	corehttp "grovehub/libs/core/http"
	"grovehub/services/integration/client/mail"
	integrationerrors "grovehub/services/integration/internal/app/constants/error"
)

func SearchHandler(client mail.Client) http.HandlerFunc {
	return func(w http.ResponseWriter,r *http.Request){
		w.Header().Set("Cache-Control","no-store")
		claims,ok:=coregateway.WorkloadFrom(r.Context())
		if !ok||claims.Subject!=mail.ReconciliationSubject||len(claims.Scopes)!=1||claims.Scopes[0]!=mail.ReconciliationScope||claims.ActorGID!=""||claims.TenantID!=""||claims.SessionID!=""{
			_ = corehttp.WriteErrorContext(r.Context(),w,http.StatusForbidden,integrationerrors.Forbidden)
			return
		}
		var in mail.SearchRequest
		r.Body=http.MaxBytesReader(w,r.Body,4096)
		if err:=corehttp.DecodeStrictJSON(r,&in);err!=nil{
			_ = corehttp.WriteErrorContext(r.Context(),w,http.StatusBadRequest,integrationerrors.Invalid)
			return
		}
		ctx,cancel:=context.WithTimeout(r.Context(),10*time.Second);defer cancel()
		if client==nil{w.WriteHeader(http.StatusServiceUnavailable);return}
		report,err:=client.Search(ctx,in)
		if errors.Is(err,mail.ErrInvalidScope){_ = corehttp.WriteErrorContext(ctx,w,http.StatusBadRequest,integrationerrors.Invalid);return}
		if err!=nil{w.WriteHeader(http.StatusServiceUnavailable);return}
		_ = corehttp.WriteSuccess(w,http.StatusOK,report)
	}
}
