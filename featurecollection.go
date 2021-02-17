package index

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func init() {
	ctx := context.Background()
	index.Register(ctx, "featurecollection", NewFeatureCollectionIndexer)
}

type FeatureCollectionIndexer struct {
	Indexer
}

func NewFeatureCollectionIndexer() (Indexer, error) {
	i := &FeatureCollectionIndexer{}
	return i, nil
}

func (i *FeatureCollectionIndexer) IndexURI(ctx context.Context, index_cb index.IndexerCallbackFunc, uri string) error {

	fh, err := ReaderWithPath(ctx, uri)

	if err != nil {
		return err
	}

	defer fh.Close()

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return err
	}

	type FC struct {
		Type     string
		Features []interface{}
	}

	var collection FC

	err = json.Unmarshal(body, &collection)

	if err != nil {
		return err
	}

	for i, f := range collection.Features {

		select {
		case <-ctx.Done():
			break
		default:
			// pass
		}

		feature, err := json.Marshal(f)

		if err != nil {
			return err
		}

		fh := bytes.NewBuffer(feature)

		path := fmt.Sprintf("%s#%d", uri, i)
		ctx = AssignPathContext(ctx, path)

		err = index_cb(ctx, fh)

		if err != nil {
			return err
		}
	}

	return nil
}
