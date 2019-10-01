package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-index"
	_ "github.com/whosonfirst/go-whosonfirst-index/fs"
	"io"
	"log"
	"strings"
	"sync/atomic"
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

	err = i.Index(ctx, paths...)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(count, i.Indexed)
}
