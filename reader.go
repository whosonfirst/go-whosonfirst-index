package index

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-index"
	"io"
	"os"
	_ "path/filepath"
)

func ReaderWithPath(ctx context.Context, abs_path string) (io.ReadSeekCloser, error) {

	if abs_path == index.STDIN {
		return os.Stdin, nil
	}

	fh, err := os.Open(abs_path)

	if err != nil {
		return nil, err
	}

	return fh, nil
}
