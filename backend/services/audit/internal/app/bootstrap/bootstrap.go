package bootstrap

import "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/module"

type Runtime interface { module.Module }
type Bootstrap interface { Build() (Runtime,error) }
