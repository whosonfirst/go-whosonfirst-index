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
	emitter_desc := fmt.Sprintf("A valid whosonfirst/go-whosonfirst-index/v2/emitter URI. Supported emitter URI schemes are: %s", valid_schemes)

	var emitter_uri = flag.String("emitter-uri", "repo://", emitter_desc)

	as_json := flag.Bool("json", false, "Emit features as a well-formed JSON array.")
	as_geojson := flag.Bool("geojson", false, "Emit features as a well-formed GeoJSON FeatureCollection record.")

	to_stdout := flag.Bool("stdout", true, "Publish features to STDOUT.")
	to_devnull := flag.Bool("null", false, "Publish features to /dev/null")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Publish features from one or more whosonfirst/go-whosonfirst-index/v2/emitter sources.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] uri(N) uri(N)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n\n")
		flag.PrintDefaults()
	}

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

	_, err := pub.Publish(ctx, *emitter_uri, uris...)

	if err != nil {
		log.Fatalf("Failed to emit features, %v", err)
	}

}
