package index

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	"net/url"
)

func IndexerInitializationFunc func(context.Context, string) (Indexer, error)

type Indexer interface {
	IndexURI(context.Context, IndexerFunc, string) error
}

var indexers roster.Roster

func ensureSpatialRoster() error {

	if indexers == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		indexers = r
	}

	return nil
}

func RegisterIndexer(ctx context.Context, scheme string, f IndexerInitializeFunc) error {

	err := ensureSpatialRoster()

	if err != nil {
		return err
	}

	return indexers.Register(ctx, scheme, f)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureSpatialRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range indexers.Indexers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func NewIndexer(ctx context.Context, uri string) (Indexer, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := indexers.Indexer(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(IndexerInitializeFunc)
	return f(ctx, uri)
}
