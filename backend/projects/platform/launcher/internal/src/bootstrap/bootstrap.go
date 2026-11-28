package bootstrap

type Bootstrap interface{ Build() (Runtime, error) }

type Runtime interface {
	Run() error
	Close()
}

// Run owns the built runtime and always closes it, including on server failure.
func Run(builder Bootstrap) error {
	runtime, err := builder.Build()
	if err != nil {
		return err
	}

	defer runtime.Close()
	return runtime.Run()
}
