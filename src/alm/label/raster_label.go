package label

import (
	"asset-tracker/src/core/asset"
	"bytes"
	"fmt"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/rasterizer"
	"image/png"
)

func RenderRaster(a *asset.Asset, params Params) (data []byte, err error) {
	svg := bytes.Buffer{}
	if e := renderToWriter(&svg, a, params); e != nil {
		err = fmt.Errorf("failed to render SVG: %w", e)
		return
	}

	c, e := canvas.ParseSVG(&svg)
	if e != nil {
		err = fmt.Errorf("failed to parse generated SVG: %w", e)
		return
	}

	r := rasterizer.New(400, 300, canvas.DefaultResolution, canvas.DefaultColorSpace)
	c.RenderTo(r)

	buf := bytes.Buffer{}
	if e := png.Encode(&buf, r.Image); e != nil {
		err = fmt.Errorf("cannot encode png: %w", err)
		return
	}

	return buf.Bytes(), nil
}
