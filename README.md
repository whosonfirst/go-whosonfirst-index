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
		log.Println("PATH", path)
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

_To be written_

## URIs and Schemes 

_To be written_

## Filters

_To be written_


## Tools

```
$> make cli
go build -mod vendor -o bin/count cmd/count/main.go
go build -mod vendor -o bin/emit cmd/emit/main.go
```

### emit

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