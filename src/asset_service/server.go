package main

import (
	"asset-tracker/src/asset_manager"
	"asset-tracker/src/core/asset"
	"asset-tracker/src/proto/asset_common"
	"asset-tracker/src/proto/asset_service"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type assetServer struct {
	Logger       *zap.Logger
	AssetManager asset_manager.AssetManager
}

const (
	NameLenMin        = 3
	NameLenMax        = 24
	DescriptionLenMin = 3
	DescriptionLenMax = 24
)

func (server *assetServer) CreateAsset(ctx context.Context, request *asset_service.CreateAssetRequest) (*asset_service.CreateAssetResponse, error) {
	if l := len(request.GetName()); l < NameLenMin || l > NameLenMax {
		return nil, status.Errorf(codes.InvalidArgument, "Asset name must contain between %v and %v characters.",
			NameLenMin, NameLenMax)
	}
	if l := len(request.GetDescription()); l < DescriptionLenMin || l > DescriptionLenMax {
		return nil, status.Errorf(codes.InvalidArgument, "Asset description contain between %v and %v characters.",
			DescriptionLenMin, DescriptionLenMax)
	}

	// TODO: Include the custom properties there.
	a := asset.New(request.GetName(), request.GetDescription())
	if err := server.AssetManager.CreateAsset(a); err != nil {
		assetIdBytes, _ := a.Id.MarshalBinary()
		server.Logger.Error(
			"Failed to write asset to datastore.",
			zap.ByteString("assetId", assetIdBytes),
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal, "Internal service error.")
	}

	return &asset_service.CreateAssetResponse{
		Asset: &asset_common.AssetObject{
			Id:          asset.EncodeIdToString(a.Id),
			Name:        a.Name,
			Description: a.Description,
			Attributes:  []*asset_common.AssetAttribute{},
		},
	}, nil
}
