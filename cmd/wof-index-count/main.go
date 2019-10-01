package main

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-index"
	_ "github.com/whosonfirst/go-whosonfirst-index/driver"
	"io"
	"log"
	"sync/atomic"
)

func main() {

	var mode = flag.String("mode", "repo://", "")
	flag.Parse()

	var count int64
	count = 0

	f := func(ctx context.Context, fh io.Reader, args ...interface{}) error {

		_, err := index.PathForContext(ctx)

		if err != nil {
			return err
		}

		atomic.AddInt64(&count, 1)
		return nil
	}

	i, err := index.NewIndexer(*mode, f)

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	paths := flag.Args()

	err = i.IndexPaths(ctx, paths...)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(count, i.Indexed)
}
