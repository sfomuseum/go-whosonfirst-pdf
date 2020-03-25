# go-whosonfirst-pdf

![](docs/images/sfom-collection.png)

Archive one or more Who's On First documents in a PDF file using an OCR-compatible font.

## Important

This is work in progress. It should be considered to work until it doesn't. Patches and other contributions are welcome.

All of the heavy-lifting is done by the [go-archive-pdf](https://github.com/sfomuseum/go-archive-pdf) package.

## Example

_Error handling omitted

### cmd/book/main.go

```
package main

import (
	_ "github.com/whosonfirst/go-whosonfirst-index/fs"
)

import (
	"context"
	"flag"
	"github.com/sfomuseum/go-archive-pdf"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-index"
	"io"
)

func main() {

	mode := flag.String("mode", "repo", "...")

	flag.Parse()

	ctx := context.Background()

	opts := pdf.NewDefaultBookOptions()
	bk, _ := pdf.NewBook(opts)

	cb := func(ctx context.Context, fh io.Reader, args ...interface{}) error {
		f, _ := feature.LoadGeoJSONFeatureFromReader(fh)
		bk.AddRecord(ctx, f.Bytes())
		return nil
	}

	idx, _ := index.NewIndexer(*mode, cb)

	uris := flag.Args()
	idx.Index(ctx, uris...)

	bk.Save("test.pdf")
}
```

## See also

* https://github.com/sfomuseum/go-archive-pdf