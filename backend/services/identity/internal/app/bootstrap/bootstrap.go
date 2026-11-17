package bootstrap
import("grovehub/libs/core/module";coregateway "grovehub/libs/core/gateway";accountclient "grovehub/services/identity/client/account")
type Runtime interface{module.Module;coregateway.AccessVerifier;AccountClient()accountclient.Client}
type Bootstrap interface{Build()(Runtime,error)}
