package index

import (
	"context"
)

func init() {
	ctx := context.Background()
	RegisterIndexer(ctx, "file", NewFileIndexer)
}

type FileIndexer struct {
	Indexer
}

func NewFileIndexer(ctx context.Context, uri string) (Indexer, error) {
	i := &FileIndexer{}
	return i, nil
}

func (i *FileIndexer) IndexURI(ctx context.Context, index_cb IndexerCallbackFunc, uri string) error {

	fh, err := ReaderWithPath(uri)

	if err != nil {
		return err
	}

	defer fh.Close()

	ctx = AssignPathContext(ctx, uri)
	return index_cb(ctx, fh)
}
