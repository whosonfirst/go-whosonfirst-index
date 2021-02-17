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
	filters *Filters
}

func NewFileIndexer(ctx context.Context, uri string) (Indexer, error) {

	f, err := NewFiltersFromURI(ctx, uri)

	if err != nil {
		return nil, err
	}

	idx := &FileIndexer{
		filters: f,
	}

	return idx, nil
}

func (idx *FileIndexer) IndexURI(ctx context.Context, index_cb IndexerCallbackFunc, uri string) error {

	fh, err := ReaderWithPath(ctx, uri)

	if err != nil {
		return err
	}

	defer fh.Close()

	if idx.filters != nil {

		ok, err := idx.filters.Apply(ctx, fh)

		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		_, err = fh.Seek(0, 0)

		if err != nil {
			return err
		}
	}

	ctx = AssignPathContext(ctx, uri)
	return index_cb(ctx, fh)
}
