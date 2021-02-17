package indexer

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-index/v2/emitter"
	"io"
	"log"
	"net/url"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

type Indexer struct {
	Emitter             emitter.Emitter
	EmitterCallbackFunc emitter.EmitterCallbackFunc
	Logger              *log.Logger
	Indexed             int64
	count               int64
	max_procs           int
}

func NewIndexer(ctx context.Context, uri string, cb emitter.EmitterCallbackFunc) (*Indexer, error) {

	idx, err := emitter.NewEmitter(ctx, uri)

	if err != nil {
		return nil, err
	}

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	max_procs := runtime.NumCPU()

	if q.Get("_max_procs") != "" {

		max, err := strconv.ParseInt(q.Get("_max_procs"), 10, 64)

		if err != nil {
			return nil, err
		}

		max_procs = int(max)
	}

	logger := log.Default()

	i := Indexer{
		Emitter:             idx,
		EmitterCallbackFunc: cb,
		Logger:              logger,
		Indexed:             0,
		count:               0,
		max_procs:           max_procs,
	}

	return &i, nil
}

func (idx *Indexer) Index(ctx context.Context, uris ...string) error {

	t1 := time.Now()

	defer func() {
		t2 := time.Since(t1)
		idx.Logger.Printf("time to index paths (%d) %v", len(uris), t2)
	}()

	idx.increment()
	defer idx.decrement()

	counter_func := func(ctx context.Context, fh io.ReadSeekCloser, args ...interface{}) error {

		defer atomic.AddInt64(&idx.Indexed, 1)
		defer fh.Close()

		return idx.EmitterCallbackFunc(ctx, fh, args...)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	procs := idx.max_procs
	throttle := make(chan bool, procs)

	for i := 0; i < procs; i++ {
		throttle <- true
	}

	done_ch := make(chan bool)
	err_ch := make(chan error)

	remaining := len(uris)

	for _, uri := range uris {

		go func(uri string) {

			<-throttle

			defer func() {
				throttle <- true
				done_ch <- true
			}()

			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}

			err := idx.Emitter.IndexURI(ctx, counter_func, uri)

			if err != nil {
				err_ch <- err
			}
		}(uri)
	}

	for remaining > 0 {
		select {
		case <-done_ch:
			remaining -= 1
		case err := <-err_ch:
			return err
		default:
			// pass
		}
	}

	return nil
}

func (idx *Indexer) IsIndexing() bool {

	if atomic.LoadInt64(&idx.count) > 0 {
		return true
	}

	return false
}

func (idx *Indexer) increment() {
	atomic.AddInt64(&idx.count, 1)
}

func (idx *Indexer) decrement() {
	atomic.AddInt64(&idx.count, -1)
}
