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
	NameLenMax        = 16
	DescriptionLenMin = 3
	DescriptionLenMax = 18
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
	if request.GetAssetKind() == asset_common.AssetKind_ASSET_KIND_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "Asset kind must be specified.")
	}

	var containerID *asset.Id

	if request.ContainerId != nil {
		if parsed, err := asset.ParseId(request.GetContainerId()); err != nil {
			return nil, status.Error(codes.InvalidArgument, "Container ID is invalid.")
		} else {
			containerID = &parsed
		}

		container, err := s.AssetManager.GetAsset(containerID)
		if err != nil {
			if errors.Is(err, asset_manager.ErrAssetNotFound) {
				return nil, status.Error(codes.NotFound, "Container does not exist.")
			}

			s.Logger.Error("Failed verifying that container asset exists.", zap.Error(err))
			return nil, status.Error(codes.Internal, "Internal error.")
		}

		if container.Kind != asset.KindContainer {
			return nil, status.Errorf(codes.InvalidArgument, "Asset %v cannot be used as container.",
				request.GetContainerId())
		}
	}

	var a asset.Asset
	if request.GetAssetKind() == asset.KindContainer {
		a = asset.NewContainer(request.GetName(), request.GetDescription(), containerID)
	} else {
		a = asset.NewItem(request.GetName(), request.GetDescription(), containerID)
	}

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
			AssetKind:   request.GetAssetKind(),
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
			ContainerId: containerIdToString(a.ContainerId),
			Name:        a.Name,
			Description: a.Description,
			AssetKind:   assetKindToGRPC(a.Kind),
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
			ContainerId: containerIdToString(a.ContainerId),
			Name:        a.Name,
			Description: a.Description,
			AssetKind:   assetKindToGRPC(a.Kind),
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

func (s *assetServer) ListAssetsInContainer(
	ctx context.Context,
	request *asset_service.ListAssetsInContainerRequest,
) (*asset_service.ListAssetsInContainerResponse, error) {
	containerID, err := asset.ParseId(request.GetContainerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid container ID.")
	}

	var nextToken string
	if request.NextToken != nil {
		nextToken = *request.NextToken
	}

	r, err := s.AssetManager.ListAssetsInContainer(&asset_manager.ListAssetsInContainerParams{
		ContainerId:  containerID,
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
			ContainerId: containerIdToString(a.ContainerId),
			Name:        a.Name,
			Description: a.Description,
			AssetKind:   assetKindToGRPC(a.Kind),
			// TODO: Handle attributes.
			Attributes: nil,
		}
	}

	var outputNextToken *string
	if r.HasNextToken {
		outputNextToken = proto.String(r.NextToken)
	}

	return &asset_service.ListAssetsInContainerResponse{
		Assets:    assets,
		NextToken: outputNextToken,
	}, nil
}
