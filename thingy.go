package index

import (
	"context"
	"io"
	"sync/atomic"
	"time"
)

const (
	STDIN = "STDIN"
)

type ThingyContextKey string

type Thingy struct {
	Indexer Indexer
	Func    IndexerCallbackFunc
	Logger  *log.Logger
	Indexed int64
	count   int64
}

func NewThingy(ctx context.Context, uri string, cb IndexerCallbackFunc) (*Thingy, error) {

	idx, err := NewIndexer(ctx, uri)

	if err != nil {
		return nil, err
	}

	logger := log.Default()

	i := Thingy{
		Indexer: idx,
		Func:    cb,
		Logger:  logger,
		Indexed: 0,
		count:   0,
	}

	return &i, nil
}

func (i *Thingy) Index(ctx context.Context, paths ...string) error {

	t1 := time.Now()

	defer func() {
		t2 := time.Since(t1)
		i.Logger.Status("time to index paths (%d) %v", len(paths), t2)
	}()

	i.increment()
	defer i.decrement()

	counter_func := func(ctx context.Context, fh io.Reader, args ...interface{}) error {
		defer atomic.AddInt64(&i.Indexed, 1)
		return i.Func(ctx, fh, args...)
	}

	for _, path := range paths {

		select {
		case <-ctx.Done():
			break
		default:
			// pass
		}

		err := i.Indexer.IndexURI(ctx, counter_func, path)

		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Thingy) IsIndexing() bool {

	if atomic.LoadInt64(&i.count) > 0 {
		return true
	}

	return false
}

func (i *Thingy) increment() {
	atomic.AddInt64(&i.count, 1)
}

func (i *Thingy) decrement() {
	atomic.AddInt64(&i.count, -1)
}
