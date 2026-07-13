package health

import "context"

// Pinger is the interface for components that can report their health status.
type Pinger interface {
	Ping(ctx context.Context) error
}
