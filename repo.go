package index

import (
	"context"
	"path/filepath"
)

func init() {
	ctx := context.Background()
	Register(ctx, "repo", NewRepoIndexer)
}

type RepoIndexer struct {
	Indexer
	indexer Indexer
}

func NewRepoIndexer(ctx context.Context, uri string) (Indexer, error) {

	directory_idx, err := NewDirectoryIndexer(ctx, uri)

	if err != nil {
		return nil, err
	}
	
	idx := &RepoIndexer{
		indexer: directory_idx,
	}

	return dr
}

func (idx *RepoIndexer) IndexURI(ctx context.Context, index_cb IndexerCallbackFunc, uri string) error {

	abs_path, err := filepath.Abs(uri)

	if err != nil {
		return err
	}

	data_path := filepath.Join(abs_path, "data")

	return idx.directory_indexer.IndexURI(ctx, index_cb, data_path)
}
