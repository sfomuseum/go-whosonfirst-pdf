package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {

	sep := flag.String("record-separator", "RECORDSEPARATOR", "...")

	flag.Parse()

	for _, path := range flag.Args() {

		fh, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open '%s', %v", path, err)
		}

		defer fh.Close()

		body, err := ioutil.ReadAll(fh)

		if err != nil {
			log.Fatalf("Failed to read '%s', %v", path, err)
		}

		str_body := string(body)

		parts := strings.Split(str_body, *sep)

		for _, raw := range parts {

			raw = strings.TrimSpace(raw)
			raw = strings.Replace(raw, "\n", "", -1)

			if !strings.HasPrefix(raw, "{") || !strings.HasSuffix(raw, "}") {
				continue
			}

			r := strings.NewReader(raw)

			f, err := feature.LoadGeoJSONFeatureFromReader(r)

			if err != nil {
				log.Printf("Failed to parse record '%s', %v", raw, err)
				continue
			}

			log.Println(f.Id())
		}
	}
}
