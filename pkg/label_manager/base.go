package label_manager

import (
	"asset-tracker/pkg/core/asset"
	"asset-tracker/pkg/label"
	"errors"
)

var ErrLabelNotFound = errors.New("asset not found")

type LabelManager interface {
	CreateLabel(params *CreateLabelParams) (*label.Label, error)
	GetLabel(id label.Id) (*label.Label, error)
	ListLabelsForAsset(assetId asset.Id) ([]label.Id, error)
	GetRenderedImageURL(id label.Id) (string, error)
}

type CreateLabelParams struct {
	AssetId               asset.Id
	FirstLine, SecondLine string
}
