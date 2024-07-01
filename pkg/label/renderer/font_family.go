package renderer

import (
	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/tdewolff/canvas"
	"log"
	"path"
)

const bazelRoot = "__main__"
const fontSize = 24

var courierNew = prepareCourierNew()
var boldFontFace = courierNew.Face(fontSize, canvas.FontBold)
var regularFontFace = courierNew.Face(fontSize, canvas.FontRegular)

func prepareCourierNew() *canvas.FontFamily {
	ff := canvas.NewFontFamily("CourierNew")

	fp := mustRetrieveFontPath("CourierNew.ttf")
	ff.MustLoadFontFile(fp, canvas.FontRegular)

	fp = mustRetrieveFontPath("CourierNewBold.ttf")
	ff.MustLoadFontFile(fp, canvas.FontBold)

	return ff
}

func mustRetrieveFontPath(name string) string {
	p := path.Join(bazelRoot, "pkg/label/renderer/data/fonts", name)
	abspath, err := runfiles.Rlocation(p)
	if err != nil {
		log.Panicf("could not retrieve font file: %v", err)
	}

	return abspath
}
