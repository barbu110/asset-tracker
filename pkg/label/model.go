package label

import "asset-tracker/pkg/core/asset"

type Label struct {
	Id                    Id
	AssetId               asset.Id
	FirstLine, SecondLine string
}
