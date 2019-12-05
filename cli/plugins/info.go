package plugins

import "context"

// Informer can be implemented to add more checks
// to `buffalo info`
type Informer interface {
	Info(ctx context.Context, args []string) error
}
