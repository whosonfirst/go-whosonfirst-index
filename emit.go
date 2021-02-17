package index

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
)

type Emitter struct {
	AsJSON    bool
	AsGeoJSON bool
	Writer    io.Writer
}

func (e *Emitter) Emit(ctx context.Context, indexer_uri string, uris ...string) (int64, error) {

	mu := new(sync.RWMutex)

	var count int64
	var count_bytes int64

	count = 0
	count_bytes = 0

	if e.AsGeoJSON {

		b, err := e.Writer.Write([]byte(`{"type":"FeatureCollection", "features":`))

		if err != nil {
			return atomic.LoadInt64(&count_bytes), err
		}

		atomic.AddInt64(&count_bytes, int64(b))
	}

	if e.AsGeoJSON || e.AsJSON {

		b, err := e.Writer.Write([]byte(`[`))

		if err != nil {
			return atomic.LoadInt64(&count_bytes), err
		}

		atomic.AddInt64(&count_bytes, int64(b))
	}

	cb := func(ctx context.Context, fh io.ReadSeekCloser, args ...interface{}) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		mu.Lock()
		defer mu.Unlock()

		atomic.AddInt64(&count, 1)

		if e.AsGeoJSON || e.AsJSON {
			if atomic.LoadInt64(&count) > 1 {

				b, err := e.Writer.Write([]byte(`,`))

				if err != nil {
					return err
				}

				atomic.AddInt64(&count_bytes, int64(b))
			}
		}

		b, err := io.Copy(e.Writer, fh)

		if err != nil {
			return err
		}

		atomic.AddInt64(&count_bytes, int64(b))
		return nil
	}

	idx, err := NewThingy(ctx, indexer_uri, cb)

	if err != nil {
		return atomic.LoadInt64(&count_bytes), err
	}

	err = idx.Index(ctx, uris...)

	if err != nil {
		return atomic.LoadInt64(&count_bytes), err
	}

	if e.AsGeoJSON || e.AsJSON {

		b, err := e.Writer.Write([]byte(`]`))

		if err != nil {
			return atomic.LoadInt64(&count_bytes), err
		}

		atomic.AddInt64(&count_bytes, int64(b))
	}

	if e.AsGeoJSON {

		b, err := e.Writer.Write([]byte(`}`))

		if err != nil {
			return atomic.LoadInt64(&count_bytes), err
		}

		atomic.AddInt64(&count_bytes, int64(b))
	}

	return atomic.LoadInt64(&count_bytes), nil
}
