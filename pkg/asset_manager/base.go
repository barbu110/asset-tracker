package asset_manager

import (
	"asset-tracker/pkg/core/asset"
	"errors"
)

var (
	ErrAssetNotFound = errAssetNotFound()
)

type AssetManager interface {
	CreateAsset(asset asset.Asset) error
	GetAsset(id *asset.Id) (a *asset.Asset, err error)
	ListAssets(params *ListAssetsParams) (data PaginatedData[asset.Asset], err error)
}

type ListAssetsParams struct {
	MaxItems     uint64
	NextToken    string
	HasNextToken bool
}

func errAssetNotFound() error {
	return errors.New("asset not found")
}
