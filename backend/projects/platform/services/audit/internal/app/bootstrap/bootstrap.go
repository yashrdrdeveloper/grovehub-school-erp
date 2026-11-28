package bootstrap

import "grovehub/libs/core/module"

type Runtime interface {
	module.Module
	Ready() bool
}
type Bootstrap interface{ Build() (Runtime, error) }
