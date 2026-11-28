package server
import "context"
type Server interface { Run(context.Context) error }
