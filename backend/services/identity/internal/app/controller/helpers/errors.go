package helpers
import(
 "errors"
 "strings"
 coreerrors "grovehub/libs/core/errors"
 definitions "grovehub/services/identity/internal/app/constants/error"
 accountrepo "grovehub/services/identity/internal/app/repository/account"
 sessionrepo "grovehub/services/identity/internal/app/repository/session"
 tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
)
// Map ensures internal or vendor errors are never serialized and keeps one stable envelope.
func Map(err error)error{
 if err==nil{return nil}
 var app *coreerrors.Error;if errors.As(err,&app){return err}
 switch{
 case errors.Is(err,accountrepo.ErrNotFound):return coreerrors.NewError(definitions.NotFound)
 case errors.Is(err,accountrepo.ErrConflict),errors.Is(err,accountrepo.ErrVersionConflict):return coreerrors.NewErrorWithMetadata(definitions.Conflict,map[string]any{"resource":"account"})
 case errors.Is(err,tokenrepo.ErrAttemptsExceeded):return coreerrors.NewError(definitions.SetupAttemptsExceeded)
 case errors.Is(err,sessionrepo.ErrInvalidSession),errors.Is(err,tokenrepo.ErrInvalidToken):return coreerrors.NewError(definitions.InvalidToken)
 }
 t:=err.Error()
 switch{
 case strings.Contains(t,"invalid credentials"),strings.Contains(t,"invalid refresh token"):return coreerrors.NewError(definitions.InvalidCredentials)
 case strings.Contains(t,"invalid proof"),strings.Contains(t,"invalid reset grant"),strings.Contains(t,"invalid or expired"):return coreerrors.NewError(definitions.InvalidToken)
 case strings.Contains(t,"authentication required"):return coreerrors.NewError(definitions.Unauthorized)
 case strings.Contains(t,"forbidden"),strings.Contains(t,"unauthorized tenant"),strings.Contains(t,"recent authentication"):return coreerrors.NewError(definitions.Forbidden)
 case strings.Contains(t,"already"),strings.Contains(t,"version conflict"),strings.Contains(t,"status transition"):return coreerrors.NewError(definitions.Conflict)
 case strings.Contains(t,"username"):return coreerrors.NewErrorWithMetadata(definitions.InvalidRequest,map[string]any{"field":"username"})
 case strings.Contains(t,"invalid"),strings.Contains(t,"must be"),strings.Contains(t,"required"):return coreerrors.NewError(definitions.InvalidRequest)
 default:return coreerrors.NewError(definitions.Internal)
 }
}
