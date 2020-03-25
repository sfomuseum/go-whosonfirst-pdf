# go-whosonfirst-pdf

![](docs/images/sfom-collection.png)

Archive one or more Who's On First documents in a PDF file using an OCR-compatible font.

## Important

This is work in progress. It should be considered to work until it doesn't. Patches and other contributions are welcome.

It's also very possible that this will be split in to two separate packages: One to deal with generating OCR-friendly PDF files for arbitrary JSON (or really any "encodable" document) and another specific to Who's On First documents. It's early days still.

## Example

_Error handling omitted

### Simple

```
package main

import (
	"context"
	"github.com/sfomuseum/go-whosonfirst-pdf"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"os"
)

func main() {

	ctx := context.Background()

	fh, _ := os.Open("example.geojson")
	f, _ := feature.LoadGeoJSONFeatureFromReader(fh)
		
	opts := pdf.NewDefaultBookOptions()
	bk, _ := pdf.NewBook(opts)

	bk.AddFeature(ctx, f)
	bk.Save("test.pdf")
}
```

### Fancy

```
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
)

func main() {

	mode := flag.String("mode", "repo", "...")

	flag.Parse()

	ctx := context.Background()

	opts := pdf.NewDefaultBookOptions()
	bk, _ := pdf.NewBook(opts)

	cb := func(ctx context.Context, fh io.Reader, args ...interface{}) error {
		f, _ := feature.LoadGeoJSONFeatureFromReader(fh)
		bk.AddFeature(ctx, f)
		return nil
	}

	idx, _ := index.NewIndexer(*mode, cb)

	uris := flag.Args()
	idx.Index(ctx, uris...)

	bk.Save("test.pdf")
}
```

## See also

* https://github.com/jung-kurt/gofpdf
* https://github.com/sfomuseum/go-font-ocra