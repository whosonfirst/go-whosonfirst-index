package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-index/v2/emitter"
	"github.com/whosonfirst/go-whosonfirst-index/v2/publisher"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	valid_schemes := strings.Join(emitter.Schemes(), ",")
	dsn_desc := fmt.Sprintf("Valid URI schemes are: %s", valid_schemes)

	indexer_uri := flag.String("indexer-uri", "repo://", dsn_desc)

	as_json := flag.Bool("json", false, "...")
	as_geojson := flag.Bool("geojson", false, "...")

	to_stdout := flag.Bool("stdout", true, "...")
	to_devnull := flag.Bool("null", false, "...")

	flag.Parse()

	if *as_geojson {
		*as_json = true
	}

	ctx := context.Background()

	writers := make([]io.Writer, 0)

	if *to_stdout {
		writers = append(writers, os.Stdout)
	}

	if *to_devnull {
		writers = append(writers, io.Discard)
	}

	wr := io.MultiWriter(writers...)

	pub := &publisher.FeaturePublisher{
		AsJSON:    *as_json,
		AsGeoJSON: *as_geojson,
		Writer:    wr,
	}

	uris := flag.Args()

	_, err := pub.Publish(ctx, *indexer_uri, uris...)

	if err != nil {
		log.Fatalf("Failed to emit features, %v", err)
	}

}
