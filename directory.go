package index

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"os"
	"path/filepath"
)

func init() {
	ctx := context.Background()
	RegisterIndexer(ctx, "directory", NewDirectoryIndexer)
}

type DirectoryIndexer struct {
	Indexer
}

func NewDirectoryIndexer(ctx context.Context, uri string) (Indexer, error) {
	i := &DirectoryIndexer{}
	return i, nil
}

func (i *DirectoryIndexer) IndexURI(ctx context.Context, index_cb IndexerCallbackFunc, uri string) error {

	abs_path, err := filepath.Abs(uri)

	if err != nil {
		return err
	}

	crawl_cb := func(path string, info os.FileInfo) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		if info.IsDir() {
			return nil
		}

		fh, err := ReaderWithPath(ctx, path)

		if err != nil {
			return err
		}

		defer fh.Close()

		ctx = AssignPathContext(ctx, path)
		return index_cb(ctx, fh)
	}

	c := crawl.NewCrawler(abs_path)
	return c.Crawl(crawl_cb)
}
