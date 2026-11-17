package bootstrap
import("grovehub/libs/core/module";accountclient "grovehub/services/identity/client/account")
type Runtime interface{module.Module;AccountClient()accountclient.Client}
type Bootstrap interface{Build()(Runtime,error)}
