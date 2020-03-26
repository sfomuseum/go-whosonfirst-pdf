package pdf

import (
	"context"
	"encoding/json"
	"github.com/jung-kurt/gofpdf"
	// "github.com/rainycape/unidecode"
	"github.com/sfomuseum/go-font-ocra"
	"log"
	"strings"
	"sync"
)

type BookOptions struct {
	Orientation     string
	Size            string
	Width           float64
	Height          float64
	DPI             float64
	Border          float64
	FontSize        float64
	Debug           bool
	OCRA            bool
	RecordSeparator string
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

type Book struct {
	PDF     *gofpdf.Fpdf
	Mutex   *sync.Mutex
	Border  BookBorder
	Canvas  BookCanvas
	Options *BookOptions
	pages   int
}

func NewDefaultBookOptions() *BookOptions {

	opts := &BookOptions{
		Orientation:     "P",
		Size:            "letter",
		Width:           0.0,
		Height:          0.0,
		DPI:             150.0,
		Border:          0.01,
		Debug:           false,
		FontSize:        12.0,
		RecordSeparator: "RECORDSEPARATOR",
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

	if opts.OCRA {

		font, err := ocra.LoadFPDFFont()

		if err != nil {
			return nil, err
		}

		pdf.AddFontFromBytes(font.Family, font.Style, font.JSON, font.Z)
		pdf.SetFont(font.Family, "", opts.FontSize)

	} else {
		pdf.SetFont("Courier", "", opts.FontSize)
	}

	w, h, _ := pdf.PageSize(1)

	page_w := w * opts.DPI
	page_h := h * opts.DPI

	border_top := 1.0 * opts.DPI
	border_bottom := border_top * 1.5
	border_left := border_top * 1.0
	border_right := border_top * 1.0

	canvas_w := page_w - (border_left + border_right)
	canvas_h := page_h - (border_top + border_bottom)

	pdf.AddPage()

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

	mu := new(sync.Mutex)

	pb := Book{
		PDF:     pdf,
		Mutex:   mu,
		Border:  b,
		Canvas:  c,
		Options: opts,
		pages:   0,
	}

	return &pb, nil
}

func (bk *Book) AddRecord(ctx context.Context, body []byte) error {

	var stub interface{}
	err := json.Unmarshal(body, &stub)

	if err != nil {
		return err
	}

	enc, err := json.Marshal(stub)

	if err != nil {
		return err
	}

	str_body := string(enc)
	str_body = strings.Replace(str_body, "\n", "", -1)

	bk.Mutex.Lock()
	defer bk.Mutex.Unlock()

	_, lh := bk.PDF.GetFontSize()
	lh = lh * 1.3

	bk.PDF.MultiCell(0, lh, str_body, "", "left", false)
	bk.PDF.MultiCell(0, lh, bk.Options.RecordSeparator, "", "", false)
	return nil
}

func (bk *Book) Save(path string) error {

	if bk.Options.Debug {
		log.Printf("save %s\n", path)
	}

	return bk.PDF.OutputFileAndClose(path)
}
