package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-index"
	"io"
	"log"
	"runtime"
)

func main() {

	var mode = flag.String("mode", "repo", "")
	var procs = flag.Int("processes", runtime.NumCPU()*2, "")

	flag.Parse()

	runtime.GOMAXPROCS(*procs)
	count := 0

	f := func(fh io.Reader, args ...interface{}) error {

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

	log.Println(count, i.Indexed)
}
