package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-index/v2"
	"io"
	"log"
	"strings"
	"sync/atomic"
	"time"
)

func main() {

	valid_schemes := strings.Join(index.Drivers(), ",")
	dsn_desc := fmt.Sprintf("Valid DSN schemes are: %s", valid_schemes)

	var dsn = flag.String("dsn", "repo://", dsn_desc)
	flag.Parse()

	var count int64
	count = 0

	cb := func(ctx context.Context, fh io.Reader, args ...interface{}) error {

		_, err := index.PathForContext(ctx)

		if err != nil {
			return err
		}

		atomic.AddInt64(&count, 1)
		return nil
	}

	i, err := index.NewIndexer(*dsn, cb)

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	paths := flag.Args()

	t1 := time.Now()

	err = i.Index(ctx, paths...)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Counted %d records (%d) in %v\n", count, i.Indexed, time.Since(t1))
}
