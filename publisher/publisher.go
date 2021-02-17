package publisher

import (
	"context"
)

type Publisher interface {
	Publish(context.Context, string, ...string) (int64, error)
}
