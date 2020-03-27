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

## Work in progress

First generate a `test.pdf` file:

```
go run -mod vendor cmd/book/main.go /usr/local/data/sfomuseum-data-collection-classifications/
```

Then use `pdfbox` to pull out images. This is done to mimic scanning the pages of the printed PDF file:

```
java -jar pdfbox-app-2.0.18.jar PDFToImage -dpi 300 test.pdf
```

Parse the text of a sample page using `tesseract`:

```
 tesseract -c tessedit_pageseg_mode=4 -c user_defined_dpi=300 -c tessedit_write_images=true -c load_system_dawg=false -c load_freq_dawg=false -c tessedit_char_whitelist='abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789\+=' test1.jpg out
```

Then try to parse what `tesseract` thinks that text is:

```
> go run -mod vendor cmd/parse/main.go out.txt
2020/03/27 13:01:33 B64 illegal base64 data at input byte 3404 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
```

## See also

* https://github.com/sfomuseum/go-archive-pdf