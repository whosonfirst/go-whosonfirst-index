package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-index"
	"log"
)

func main() {

	var mode = flag.String("mode", "repo", "")

	flag.Parse()

	count := 0

	f := func(path string, args ...interface{}) error {
		count += 1
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

	log.Println(count)
}
