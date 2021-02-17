# go-whosonfirst-index

Go package for indexing Who's On First documents

## Example

```
package main

import (
       "context"
       "flag"
       "github.com/whosonfirst/go-whosonfirst-index/v2"
       "io"
       "log"
)

func main() {

	var uri = flag.String("indexer-uri", "repo://", "A valid go-whosonfirst-index URI")
	
     	flag.Parse()
	
	cb := func(ctx context.Context, fh io.ReadSeekCloser, args ...interface{}) error {

		path, _ := index.PathForContext(ctx)

		log.Println("PATH", path)
		return nil
	}

	i, _ := index.NewIndexer(*uri, cb)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	paths := flag.Args()

	i.Index(ctx, paths...)
}	
```

_Error handling removed for the sake of brevity._

## Tools

```
$> make cli
go build -mod vendor -o bin/count cmd/count/main.go
go build -mod vendor -o bin/emit cmd/emit/main.go
```

### emit

```
$> ./bin/emit \
	-indexer-uri 'repo://?include=properties.sfomuseum:placetype=museum' \
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

## Schemes

_To be written_