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
}

func NewFileListIndexer(ctx context.Context, uri string) (Indexer, error) {
	i := &FileListIndexer{}
	return i, nil
}

func (i *FileListIndexer) IndexURI(ctx context.Context, index_cb IndexerCallbackFunc, uri string) error {

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
