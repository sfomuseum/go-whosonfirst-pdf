package pdf

import (
	"context"
	"github.com/jung-kurt/gofpdf"
	// "github.com/rainycape/unidecode"
	"github.com/sfomuseum/go-font-ocra"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"sync"
	"log"
)

type BookOptions struct {
	Orientation string
	Size        string
	Width       float64
	Height      float64
	DPI         float64
	Border      float64
	Debug       bool
}

type BookBorder struct {
	Top    float64
	Bottom float64
	Left   float64
	Right  float64
}

type BookCanvas struct {
	Width  float64
	Height float64
}

type BookText struct {
	Font   string
	Style  string
	Size   float64
	Margin float64
	Colour []int
}

type Book struct {
	PDF      *gofpdf.Fpdf
	Mutex    *sync.Mutex
	Border   BookBorder
	Canvas   BookCanvas
	Text     BookText
	Options  *BookOptions
	pages    int
	tmpfiles []string
}

func NewDefaultBookOptions() *BookOptions {

	opts := &BookOptions{
		Orientation: "P",
		Size:        "letter",
		Width:       0.0,
		Height:      0.0,
		DPI:         150.0,
		Border:      0.01,
		Debug:       false,
	}

	return opts
}

func NewBook(opts *BookOptions) (*Book, error) {

	var pdf *gofpdf.Fpdf

	if opts.Size == "custom" {

		sz := gofpdf.SizeType{
			Wd: opts.Width,
			Ht: opts.Height,
		}

		init := gofpdf.InitType{
			OrientationStr: opts.Orientation,
			UnitStr:        "in",
			SizeStr:        "",
			Size:           sz,
			FontDirStr:     "",
		}

		pdf = gofpdf.NewCustom(&init)

	} else {

		pdf = gofpdf.New(opts.Orientation, "in", opts.Size, "")
	}

	t := BookText{
		Font:   "Helvetica",
		Style:  "",
		Size:   8.0,
		Margin: 0.1,
		Colour: []int{128, 128, 128},
	}

	// pdf.SetFont(t.Font, t.Style, t.Size)

	font, err := ocra.LoadFPDFFont()

	if err != nil {
		return nil, err
	}

	pdf.AddFontFromBytes(font.Family, font.Style, font.JSON, font.Z)
	pdf.SetFont(font.Family, "", 8.0)

	pdf.SetTextColor(t.Colour[0], t.Colour[1], t.Colour[2])

	w, h, _ := pdf.PageSize(1)

	page_w := w * opts.DPI
	page_h := h * opts.DPI

	border_top := 1.0 * opts.DPI
	border_bottom := border_top * 1.5
	border_left := border_top * 1.0
	border_right := border_top * 1.0

	canvas_w := page_w - (border_left + border_right)
	canvas_h := page_h - (border_top + border_bottom)

	pdf.SetAutoPageBreak(false, border_bottom)

	b := BookBorder{
		Top:    border_top,
		Bottom: border_bottom,
		Left:   border_left,
		Right:  border_right,
	}

	c := BookCanvas{
		Width:  canvas_w,
		Height: canvas_h,
	}

	tmpfiles := make([]string, 0)
	mu := new(sync.Mutex)

	pb := Book{
		PDF:      pdf,
		Mutex:    mu,
		Border:   b,
		Canvas:   c,
		Text:     t,
		Options:  opts,
		pages:    0,
		tmpfiles: tmpfiles,
	}

	return &pb, nil
}

func (bk *Book) AddFeature(ctx context.Context, f geojson.Feature) error {

	log.Println("ADD", f.Id())
	return nil
}
