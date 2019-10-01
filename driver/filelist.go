package driver

import (
	"bufio"
	"context"
	"github.com/whosonfirst/go-whosonfirst-index"
)

func init() {

	dr := &FileListDriver{}

	index.Register("filelist", dr)
}

type FileListDriver struct {
	index.Driver
}

func (d *FileListDriver) Open(uri string) error {
	return nil
}

func (d *FileListDriver) IndexURI(ctx context.Context, index_cb index.IndexerFunc, uri string) error {

	fh, err := readerFromPath(uri)

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

		fh, err := readerFromPath(path)

		if err != nil {
			return err
		}

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
