package renderer

import (
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/datamatrix"
	"image"
)

type renderBarcodeParams struct {
	Content string
	Size    int
}

func renderBarcode(params *renderBarcodeParams) (*image.Image, error) {
	bc, err := datamatrix.Encode(params.Content)
	if err != nil {
		return nil, fmt.Errorf("datamatrix encode failed: %w", err)
	}

	bc, err = barcode.Scale(bc, params.Size, params.Size)
	if err != nil {
		return nil, fmt.Errorf("cannot scale barcode: %w", err)
	}

	// buf := bytes.Buffer{}
	//if e := png.Encode(&buf, bc); e != nil {
	//	return nil, fmt.Errorf("cannot encode png: %w", err)
	//}
	//return dataurl.EncodeBytes(buf.Bytes()), nil

	return &bc, nil
}
