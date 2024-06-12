package main

import (
	"asset-tracker/pkg/asset_manager"
	"asset-tracker/pkg/core/asset"
	"asset-tracker/proto/asset_common"
	"asset-tracker/proto/asset_service"
	"context"
	"errors"
	"github.com/golang/protobuf/proto"
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

const (
	MsgInternalServiceError = "Internal service error."
)

func (s *assetServer) CreateAsset(ctx context.Context, request *asset_service.CreateAssetRequest) (*asset_service.CreateAssetResponse, error) {
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
	if err := s.AssetManager.CreateAsset(a); err != nil {
		assetIdBytes, _ := a.Id.MarshalBinary()
		s.Logger.Error(
			"Failed to write asset to datastore.",
			zap.ByteString("assetId", assetIdBytes),
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal, MsgInternalServiceError)
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

func (s *assetServer) GetAsset(ctx context.Context, request *asset_service.GetAssetRequest) (*asset_service.GetAssetResponse, error) {
	id, err := asset.ParseId(request.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid ID.")
	}

	a, err := s.AssetManager.GetAsset(&id)
	if errors.Is(err, asset_manager.ErrAssetNotFound) {
		return nil, status.Error(codes.NotFound, "Asset not found.")
	} else if err != nil {
		return nil, status.Error(codes.Internal, MsgInternalServiceError)
	}

	return &asset_service.GetAssetResponse{
		Asset: &asset_common.AssetObject{
			Id:          asset.EncodeIdToString(a.Id),
			Name:        a.Name,
			Description: a.Description,
			Attributes:  []*asset_common.AssetAttribute{},
		},
	}, nil
}

func (s *assetServer) ListAssets(ctx context.Context, request *asset_service.ListAssetsRequest) (*asset_service.ListAssetsResponse, error) {
	var nextToken string
	if request.NextToken != nil {
		nextToken = *request.NextToken
	}

	r, err := s.AssetManager.ListAssets(&asset_manager.ListAssetsParams{
		MaxItems:     request.MaxItems,
		NextToken:    nextToken,
		HasNextToken: request.NextToken != nil,
	})
	if err != nil {
		if errors.Is(err, asset_manager.ErrInvalidNextToken) {
			return nil, status.Error(codes.InvalidArgument, "Provided nextToken is invalid.")
		}

		s.Logger.Error("Failed to list assets.", zap.Error(err))
		return nil, status.Error(codes.Internal, MsgInternalServiceError)
	}

	assets := make([]*asset_common.AssetObject, len(r.Items))
	for i, a := range r.Items {
		assets[i] = &asset_common.AssetObject{
			Id:          asset.EncodeIdToString(a.Id),
			Name:        a.Name,
			Description: a.Description,
			// TODO: Handle attributes.
			Attributes: nil,
		}
	}

	var outputNextToken *string
	if r.HasNextToken {
		outputNextToken = proto.String(r.NextToken)
	}

	return &asset_service.ListAssetsResponse{
		Assets:    assets,
		NextToken: outputNextToken,
	}, nil
}
