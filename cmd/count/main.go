package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-index/v2/emitter"
	"github.com/whosonfirst/go-whosonfirst-index/v2/indexer"
	"io"
	"log"
	"strings"
	"sync/atomic"
	"time"
)

func main() {

	valid_schemes := strings.Join(emitter.Schemes(), ",")
	dsn_desc := fmt.Sprintf("Valid DSN schemes are: %s", valid_schemes)

	var indexer_uri = flag.String("indexer-uri", "repo://", dsn_desc)
	flag.Parse()

	ctx := context.Background()

	var count int64
	count = 0

	cb := func(ctx context.Context, fh io.ReadSeekCloser, args ...interface{}) error {

		_, err := emitter.PathForContext(ctx)

		if err != nil {
			return err
		}

		atomic.AddInt64(&count, 1)
		return nil
	}

	i, err := indexer.NewIndexer(ctx, *indexer_uri, cb)

	if err != nil {
		log.Fatal(err)
	}

	paths := flag.Args()

	t1 := time.Now()

	err = i.Index(ctx, paths...)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Counted %d records (%d) in %v\n", count, i.Indexed, time.Since(t1))
}
