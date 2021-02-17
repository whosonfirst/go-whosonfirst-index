package index

import (
	"bufio"
	"context"
)

func init() {
	ctx := context.Background()
	RegisterIndexer(ctx, "filelist", NewFileListIndexer)
}

type FileListIndexer struct {
	Indexer
	filters *Filters
}

func NewFileListIndexer(ctx context.Context, uri string) (Indexer, error) {

	f, err := NewFiltersFromURI(ctx, uri)

	if err != nil {
		return nil, err
	}

	idx := &FileListIndexer{
		filters: f,
	}

	return idx, nil
}

func (idx *FileListIndexer) IndexURI(ctx context.Context, index_cb IndexerCallbackFunc, uri string) error {

	fh, err := ReaderWithPath(ctx, uri)

	if err != nil {
		return err
	}

	defer fh.Close()

	scanner := bufio.NewScanner(fh)

	for scanner.Scan() {

		select {
		case <-ctx.Done():
			break
		default:
			// pass
		}

		path := scanner.Text()

		fh, err := ReaderWithPath(ctx, path)

		if err != nil {
			return err
		}

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

		ctx = AssignPathContext(ctx, path)

		err = index_cb(ctx, fh)

		if err != nil {
			return err
		}
	}

	err = scanner.Err()

	if err != nil {
		return err
	}

	return nil
}
