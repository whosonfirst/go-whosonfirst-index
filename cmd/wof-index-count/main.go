package main

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-index"
	"io"
	"log"
	"sync/atomic"
)

func main() {

	var mode = flag.String("mode", "repo", "")
	flag.Parse()

	var count int64
	count = 0

	f := func(fh io.Reader, ctx context.Context, args ...interface{}) error {

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

	for _, path := range flag.Args() {

		err := i.IndexPath(path)

		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println(count, i.Indexed)
}
