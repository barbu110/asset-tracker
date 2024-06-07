package label

import (
	"asset-tracker/src/alm/asset"
	"bytes"
	"encoding/base64"
	"fmt"
	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/vincent-petithory/dataurl"
	"image/png"
)

const (
	FontFamily    = "Courier"
	FontSize      = 28
	SmallFontSize = 24
)

type Params struct {
	Width              int
	Height             int
	Padding            int
	UseWhiteBackground bool
}

func Render(a *asset.Asset, params Params) (data []byte, err error) {
	buffer := bytes.Buffer{}
	canvas := svg.New(&buffer)

	canvas.Start(params.Width, params.Height)
	if params.UseWhiteBackground {
		canvas.Rect(0, 0, params.Width, params.Height, "fill=\"#fff\"")
	}
	renderText(canvas, a.Name, boldTextParams(params.Padding, params.Padding))
	renderText(canvas, a.Description, defaultTextParams(params.Padding, params.Padding+FontSize+params.Padding/2.))

	const BarcodeSize = 160

	idData, _ := a.Id.MarshalBinary()
	idBase64 := base64.StdEncoding.EncodeToString(idData)
	idY := params.Height - params.Padding - SmallFontSize
	renderText(canvas, idBase64, smallTextParams(params.Padding, idY))

	bcY := params.Height - params.Padding - BarcodeSize - SmallFontSize - params.Padding
	bcUri, e := renderBarcode(canvas, idBase64, &renderBarcodeParams{params.Padding, bcY, BarcodeSize})
	if e != nil {
		err = fmt.Errorf("failed to render barcode: %w", e)
		return
	}

	canvas.Image(params.Padding, bcY, BarcodeSize, BarcodeSize, bcUri)

	canvas.End()

	data = buffer.Bytes()
	return
}

type renderTextParams struct {
	X, Y       int
	FontFamily string
	FontSize   int
	FontWeight string
}

func renderText(canvas *svg.SVG, t string, params renderTextParams) {
	styles := []string{
		"dominant-baseline=\"hanging\"",
		fmt.Sprintf("font-size=\"%v\"", params.FontSize),
		fmt.Sprintf("font-family=\"%v\"", params.FontFamily),
		fmt.Sprintf("font-weight=\"%v\"", params.FontWeight),
	}
	canvas.Text(params.X, params.Y, t, styles...)
}

func defaultTextParams(x, y int) renderTextParams {
	return renderTextParams{
		X:          x,
		Y:          y,
		FontFamily: FontFamily,
		FontSize:   FontSize,
		FontWeight: "normal",
	}
}

func smallTextParams(x, y int) renderTextParams {
	return renderTextParams{
		X:          x,
		Y:          y,
		FontFamily: FontFamily,
		FontSize:   SmallFontSize,
		FontWeight: "normal",
	}
}

func boldTextParams(x, y int) renderTextParams {
	return renderTextParams{
		X:          x,
		Y:          y,
		FontFamily: FontFamily,
		FontSize:   FontSize,
		FontWeight: "bold",
	}
}

type renderBarcodeParams struct {
	X, Y, Size int
}

func renderBarcode(canvas *svg.SVG, content string, params *renderBarcodeParams) (uri string, err error) {
	bc, err := datamatrix.Encode(content)
	if err != nil {
		return
	}

	bc, err = barcode.Scale(bc, params.Size, params.Size)
	if err != nil {
		err = fmt.Errorf("cannot scale barcode: %w", err)
		return
	}

	buf := bytes.Buffer{}
	if e := png.Encode(&buf, bc); e != nil {
		err = fmt.Errorf("cannot encode png: %w", err)
		return
	}

	uri = dataurl.EncodeBytes(buf.Bytes())
	fmt.Println(uri)
	return
}
