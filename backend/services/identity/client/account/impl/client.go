package impl
import(
 accountclient "grovehub/services/identity/client/account"
 service "grovehub/services/identity/internal/app/services/account"
)
// identityAccountClient is the concrete receiver used by the generated @Proxy
// methods. Methods delegate into the scoped Account service, not raw SQL.
type identityAccountClient struct{accounts service.Service}
func New(accounts service.Service)accountclient.Client{return &identityAccountClient{accounts:accounts}}
