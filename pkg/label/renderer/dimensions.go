package renderer

import (
	"github.com/tdewolff/canvas"
	"math"
)

const (
	printerDPI = 203
	inchPerMm  = 1.0 / 25.4
)

const (
	padding = 10
)

var LabelSizeMM = canvas.Size{W: 40, H: 30}
var LabelSizePX = labelSizePX(LabelSizeMM, printerDPI)

func labelSizePX(sizeMm canvas.Size, dpi float64) canvas.Size {
	convert := func(mm float64) float64 {
		px := mm * dpi * inchPerMm
		return math.Floor(px)
	}

	return canvas.Size{W: convert(sizeMm.W), H: convert(sizeMm.H)}
}
