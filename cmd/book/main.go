package main

import (
	_ "github.com/whosonfirst/go-whosonfirst-index/fs"
)

import (
	"context"
	"flag"
	"github.com/sfomuseum/go-whosonfirst-pdf"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-index"
	"io"
	"log"
)

func main() {

	mode := flag.String("mode", "repo", "...")

	flag.Parse()

	ctx := context.Background()

	opts := pdf.NewDefaultBookOptions()
	bk, err := pdf.NewBook(opts)

	if err != nil {
		log.Fatalf("Failed to create book, %v", err)
	}

	cb := func(ctx context.Context, fh io.Reader, args ...interface{}) error {

		path, err := index.PathForContext(ctx)

		if err != nil {
			return err
		}

		f, err := feature.LoadGeoJSONFeatureFromReader(fh)

		if err != nil {
			log.Printf("Failed to load feature for '%s', %v", path, err)
			return err
		}

		err = bk.AddFeature(ctx, f)

		if err != nil {
			log.Printf("Failed to add feature for '%s', %v", path, err)
			return err
		}

		return nil
	}

	idx, err := index.NewIndexer(*mode, cb)

	if err != nil {
		log.Fatalf("Failed to create indexer, %v", err)
	}

	uris := flag.Args()

	err = idx.Index(ctx, uris...)

	if err != nil {
		log.Fatal(err)
	}

}