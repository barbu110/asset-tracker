package main

import (
	"asset-tracker/pkg/core/asset"
	"asset-tracker/proto/asset_common"
	"fmt"
	"slices"
)

func assetKindToGRPC(assetKind asset.Kind) asset_common.AssetKind {
	switch assetKind {
	case asset.KindUnspecified:
		return asset_common.AssetKind_ASSET_KIND_UNSPECIFIED
	case asset.KindItem:
		return asset_common.AssetKind_ASSET_KIND_ITEM
	case asset.KindContainer:
		return asset_common.AssetKind_ASSET_KIND_CONTAINER
	}
	panic(fmt.Sprintf("unknown asset kind: %v", assetKind))
}

func containerIdToString(id asset.Id) *string {
	if slices.Equal(id, asset.RootContainerId()) {
		return nil
	}

	v := asset.EncodeIdToString(id)
	return &v
}
