package label_manager

import (
	"asset-tracker/pkg/label"
	"fmt"
	"path"
)

func renderedImagePNGKey(id label.Id) string {
	return path.Join("labels", fmt.Sprintf("%v.png", label.EncodeIdToString(id)))
}
