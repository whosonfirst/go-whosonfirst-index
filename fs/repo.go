package fs

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-index"
	"path/filepath"
)

func init() {

	dd := &DirectoryDriver{}

	dr := &RepoDriver{
		driver: dd,
	}

	index.Register("repo", dr)
}

type RepoDriver struct {
	index.Driver
	driver index.Driver
}

func (d *RepoDriver) Open(uri string) error {
	return d.driver.Open(uri)
}

func (d *RepoDriver) IndexURI(ctx context.Context, index_cb index.IndexerFunc, uri string) error {

	abs_path, err := filepath.Abs(uri)

	if err != nil {
		return err
	}

	data_path := filepath.Join(abs_path, "data")

	return d.driver.IndexURI(ctx, index_cb, data_path)
}
