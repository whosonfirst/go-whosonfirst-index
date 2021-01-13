package index

import (
	"bytes"
	"context"
	"github.com/aaronland/go-json-query"
	"io"
	"io/ioutil"
)

func NewIndexerWithQuerySet(ctx context.Context, indexer_uri string, indexer_cb IndexerFunc, qs *query.QuerySet) (*Indexer, error) {

	if qs == nil {
		return NewIndexer(indexer_uri, indexer_cb)
	}

	if len(qs.Queries) == 0 {
		return NewIndexer(indexer_uri, indexer_cb)
	}

	query_cb := func(ctx context.Context, fh io.Reader, args ...interface{}) error {

		body, err := ioutil.ReadAll(fh)

		if err != nil {
			return err
		}

		matches, err := query.Matches(ctx, qs, body)

		if err != nil {
			return err
		}

		if !matches {
			return nil
		}

		br := bytes.NewReader(body)
		return indexer_cb(ctx, br, args...)
	}

	return NewIndexer(indexer_uri, query_cb)
}
