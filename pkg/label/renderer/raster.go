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
	c := canvas.NewFromSize(LabelSizePX)
	ctx := canvas.NewContext(c)

	ctx.DrawText(padding, padding, canvas.NewTextLine(boldFontFace, params.FirstLine, canvas.Left))
	ctx.DrawText(padding, 2*padding+fontSize, canvas.NewTextLine(regularFontFace, params.SecondLine, canvas.Left))

	bc, err := renderBarcode(&renderBarcodeParams{
		Content: hex.EncodeToString(params.BarcodeData),
		Size:    100,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to render barcode: %w", err)
	}
	ctx.DrawImage(padding, 3*padding+2*fontSize, *bc, canvas.DefaultResolution)

	buf := bytes.Buffer{}
	if err := c.Write(&buf, renderers.PNG()); err != nil {
		return nil, fmt.Errorf("failed to write image: %w", err)
	}
	return buf.Bytes(), nil
}
