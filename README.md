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
	
	cb := func(ctx context.Context, fh io.Reader, args ...interface{}) error {

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
