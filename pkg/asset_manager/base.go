package asset_manager

import (
	"asset-tracker/pkg/core/asset"
	"asset-tracker/pkg/pagination"
	"errors"
)

var (
	ErrAssetNotFound    = errors.New("asset not found")
	ErrInvalidNextToken = errors.New("invalid next token")
)

const ListAssetsDefaultMaxItems uint64 = 100

type AssetManager interface {
	CreateAsset(asset asset.Asset) error
	GetAsset(id *asset.Id) (*asset.Asset, error)
	ListAssets(params *ListAssetsParams) (data pagination.PaginatedData[asset.Asset], err error)
}

type ListAssetsParams struct {
	MaxItems     uint64
	NextToken    string
	HasNextToken bool
}

func (p *ListAssetsParams) GetMaxItems() uint64 {
	if p.MaxItems == 0 {
		return ListAssetsDefaultMaxItems
	} else {
		return p.MaxItems
	}
}
