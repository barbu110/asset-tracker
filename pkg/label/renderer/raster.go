package renderer

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

type RasterRenderer struct{}

func (r *RasterRenderer) Render(params *RenderLabelParams) ([]byte, error) {
	c := canvas.NewFromSize(LabelSizeMM)
	ctx := canvas.NewContext(c)
	ctx.SetCoordSystem(canvas.CartesianIV)

	// Draw the background.
	ctx.SetFillColor(canvas.White)
	ctx.DrawPath(0, 0, canvas.MustParseSVGPath(fmt.Sprintf(
		"M 0 0 h %v v %v h -%v Z",
		LabelSizePX.W,
		LabelSizePX.H,
		LabelSizePX.W)))
	ctx.DrawText(padding, padding, canvas.NewTextBox(boldFontFace, params.FirstLine, LabelSizeMM.W, fontSize, canvas.Left, canvas.Top, 0, 1.0))
	// kinda sketchy, figure out numbers that make sense.
	ctx.DrawText(padding, 3*padding, canvas.NewTextBox(regularFontFace, params.SecondLine, LabelSizeMM.W, fontSize, canvas.Left, canvas.Top, 0, 1.0))

	bc, err := renderBarcode(&renderBarcodeParams{
		Content: hex.EncodeToString(params.BarcodeData),
		Size:    3000, // again, figure out some numbers that make sense
	})
	if err != nil {
		return nil, fmt.Errorf("failed to render barcode: %w", err)
	}
	ctx.DrawImage(padding, 6*padding, *bc, printerDPI)

	buf := bytes.Buffer{}
	if err := c.Write(&buf, renderers.PNG(canvas.DPI(printerDPI*zoomFactor))); err != nil {
		return nil, fmt.Errorf("failed to write image: %w", err)
	}
	return buf.Bytes(), nil
}
