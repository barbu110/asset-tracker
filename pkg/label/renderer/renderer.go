package renderer

type RenderLabelParams struct {
	FirstLine, SecondLine string
	BarcodeData           []byte
}

type LabelRenderer interface {
	Render(params *RenderLabelParams) ([]byte, error)
}
