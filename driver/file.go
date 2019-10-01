package driver

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-index"
)

func init() {

	dr := &FileDriver{}

	index.Register("file", dr)
}

type FileDriver struct {
	index.Driver
}

func (d *FileDriver) Open(uri string) error {
	return nil
}

func (d *FileDriver) IndexURI(ctx context.Context, index_cb index.IndexerFunc, uri string) error {

	fh, err := readerFromPath(uri)

	if err != nil {
		return err
	}

	return index_cb(ctx, fh)

}
