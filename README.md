# go-whosonfirst-index

Go package for indexing Who's On First documents

## Example

```
package main

import (
       "context"
       "flag"
       "github.com/whosonfirst/go-whosonfirst-index/v2/emitter"       
       "github.com/whosonfirst/go-whosonfirst-index/v2/indexer"
       "io"
       "log"
)

func main() {

	emitter_uri := flag.String("emitter-uri", "repo://", "A valid whosonfirst/go-whosonfirst-index/v2/emitter URI")
	
     	flag.Parse()

	ctx := context.Background()

	cb := func(ctx context.Context, fh io.ReadSeekCloser, args ...interface{}) error {
		path, _ := index.PathForContext(ctx)
		log.Printf("Indexing %s\n", path)
		return nil
	}

	idx, _ := indexer.NewIndexer(ctx, *emitter_uri, cb)

	uris := flag.Args()
	idx.Index(ctx, uris...)
}	
```

_Error handling removed for the sake of brevity._

## Concepts

### Indexer

### Emitters

_To be written_

## Interfaces

```
type EmitterInitializeFunc func(context.Context, string) (Emitter, error)

type EmitterCallbackFunc func(context.Context, io.ReadSeekCloser, ...interface{}) error

type Emitter interface {
	IndexURI(context.Context, EmitterCallbackFunc, string) error
}
```

_To be written_

## URIs and Schemes 

_To be written_

### directory://

### featurecollection://

### file://

### filelist://

### geojsonls://

### repo://

## Filters

_To be written_

## Differences from "v1"

There was never a "v1" release. The last published release before "v2" was [v0.3.4](https://github.com/whosonfirst/go-whosonfirst-index/releases/tag/v0.3.4).

* Go 1.16 or higher is required.
* The introduction of the `emitter.Emitter` interface separate from a general-purpose `indexer.Indexer` instance.
* Migrating the `index.NewIndexer` method in to the `indexer.NewIndexer` package method.
* Migrating the `index.PathForContext` method in to the `emitter.PathForContext` package method.
* Migrating the `index.Drivers` method in to the `emitter.Schemes` package method.
* The use of the `aaronland/go-roster` package to manage registered emitters.
* Changing the requirement in emitter (previously indexer) callbacks from `io.Reader` to `io.ReadSeekCloser`.
* The introduction of the `filters.Filters` interface, and corresponding emitter URI query parameters, for limiting results that are sent to emitter (previously indexer) callback functions.

## Tools

```
$> make cli
go build -mod vendor -o bin/count cmd/count/main.go
go build -mod vendor -o bin/emit cmd/emit/main.go
```

### count

Count files in one or more whosonfirst/go-whosonfirst-index/v2/emitter sources.

```
$> ./bin/count -h
Count files in one or more whosonfirst/go-whosonfirst-index/v2/emitter sources.
Usage:
	 ./bin/count [options] uri(N) uri(N)
Valid options are:

  -emitter-uri string
    	A valid whosonfirst/go-whosonfirst-index/v2/emitter URI. Supported emitter URI schemes are: directory://,featurecollection://,file://,filelist://,geojsonl://,repo:// (default "repo://")
```

For example:

```
$> ./bin/count \
	/usr/local/data/sfomuseum-data-architecture/

2021/02/17 14:07:01 time to index paths (1) 87.908997ms
2021/02/17 14:07:01 Counted 1072 records (1072) in 88.045771ms
```

Or:

```
$> ./bin/count \
	-emitter-uri 'repo://?include=properties.sfomuseum:placetype=terminal&include=properties.mz:is_current=1' \
	/usr/local/data/sfomuseum-data-architecture/
	
2021/02/17 14:09:18 time to index paths (1) 71.06355ms
2021/02/17 14:09:18 Counted 4 records (4) in 71.184227ms
```

### emit

Publish features from one or more whosonfirst/go-whosonfirst-index/v2/emitter sources.

```
$> ./bin/emit -h
Publish features from one or more whosonfirst/go-whosonfirst-index/v2/emitter sources.
Usage:
	 ./bin/emit [options] uri(N) uri(N)
Valid options are:

  -emitter-uri string
    	A valid whosonfirst/go-whosonfirst-index/v2/emitter URI. Supported emitter URI schemes are: directory://,featurecollection://,file://,filelist://,geojsonl://,repo:// (default "repo://")
  -geojson
    	Emit features as a well-formed GeoJSON FeatureCollection record.
  -json
    	Emit features as a well-formed JSON array.
  -null
    	Publish features to /dev/null
  -stdout
    	Publish features to STDOUT. (default true)
```

For example:

```
$> ./bin/emit \
	-emitter-uri 'repo://?include=properties.sfomuseum:placetype=museum' \
	-geojson \	
	/usr/local/data/sfomuseum-data-architecture/ \

| jq '.features[]["properties"]["wof:id"]'

1729813675
1477855937
1360521563
1360521569
1360521565
1360521571
1159157863
```

## See also

* https://github.com/aaronland/go-json-query
* https://github.com/aaronland/go-roster