package index

import (
	"context"
	"errors"
)

type IndexerContextKey string

func ContextForPath(path string) (context.Context, error) {

	ctx := AssignPathContext(context.Background(), path)
	return ctx, nil
}

func AssignPathContext(ctx context.Context, path string) context.Context {

	key := IndexerContextKey("path")
	return context.WithValue(ctx, key, path)
}

func PathForContext(ctx context.Context) (string, error) {

	k := IndexerContextKey("path")
	path := ctx.Value(k)

	if path == nil {
		return "", errors.New("path is not set")
	}

	return path.(string), nil
}
