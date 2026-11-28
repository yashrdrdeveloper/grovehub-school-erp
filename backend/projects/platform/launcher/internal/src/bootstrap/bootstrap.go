package bootstrap

type Bootstrap interface { Build() (Runtime,error) }
type Runtime interface { Run() error; Close() }
