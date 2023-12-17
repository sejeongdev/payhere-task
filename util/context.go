package util

import (
	"context"
	"time"
)

const (
	OwnerKey string = "uid"
)

// WithTimeout ...
func WithTimeout(ctx context.Context, timeout time.Duration) (newCtx context.Context, cancel func()) {
	return context.WithTimeout(ctx, timeout)
}
