package publisher

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-index/v2/indexer"
	"io"
	"sync"
	"sync/atomic"
)

type FeaturePublisher struct {
	Publisher
	AsJSON    bool
	AsGeoJSON bool
	Writer    io.Writer
}

func (pub *FeaturePublisher) Publish(ctx context.Context, emitter_uri string, uris ...string) (int64, error) {

	mu := new(sync.RWMutex)

	var count int64
	var count_bytes int64

	count = 0
	count_bytes = 0

	if pub.AsGeoJSON {

		b, err := pub.Writer.Write([]byte(`{"type":"FeatureCollection", "features":`))

		if err != nil {
			return atomic.LoadInt64(&count_bytes), err
		}

		atomic.AddInt64(&count_bytes, int64(b))
	}

	if pub.AsGeoJSON || pub.AsJSON {

		b, err := pub.Writer.Write([]byte(`[`))

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

		if pub.AsGeoJSON || pub.AsJSON {
			if atomic.LoadInt64(&count) > 1 {

				b, err := pub.Writer.Write([]byte(`,`))

				if err != nil {
					return err
				}

				atomic.AddInt64(&count_bytes, int64(b))
			}
		}

		b, err := io.Copy(pub.Writer, fh)

		if err != nil {
			return err
		}

		atomic.AddInt64(&count_bytes, int64(b))
		return nil
	}

	idx, err := indexer.NewIndexer(ctx, emitter_uri, cb)

	if err != nil {
		return atomic.LoadInt64(&count_bytes), err
	}

	err = idx.Index(ctx, uris...)

	if err != nil {
		return atomic.LoadInt64(&count_bytes), err
	}

	if pub.AsGeoJSON || pub.AsJSON {

		b, err := pub.Writer.Write([]byte(`]`))

		if err != nil {
			return atomic.LoadInt64(&count_bytes), err
		}

		atomic.AddInt64(&count_bytes, int64(b))
	}

	if pub.AsGeoJSON {

		b, err := pub.Writer.Write([]byte(`}`))

		if err != nil {
			return atomic.LoadInt64(&count_bytes), err
		}

		atomic.AddInt64(&count_bytes, int64(b))
	}

	return atomic.LoadInt64(&count_bytes), nil
}
