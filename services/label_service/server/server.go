package server

import (
	"asset-tracker/pkg/core/asset"
	"asset-tracker/pkg/label"
	"asset-tracker/pkg/label_manager"
	"asset-tracker/proto/label_service"
	"context"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LabelServer struct {
	Logger       *zap.Logger
	LabelManager label_manager.LabelManager
}

func (l *LabelServer) RenderLabel(ctx context.Context, request *label_service.RenderLabelRequest) (*label_service.RenderLabelResponse, error) {
	assetId, err := asset.ParseId(request.Asset.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Asset ID.")
	}
	model, err := l.LabelManager.CreateLabel(&label_manager.CreateLabelParams{
		AssetId:    assetId,
		FirstLine:  request.Asset.Name,
		SecondLine: request.Asset.Description,
	})
	if err != nil {
		l.Logger.Error("Failed creating label.", zap.String("assetId", request.Asset.Id), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Internal error.")
	}

	return &label_service.RenderLabelResponse{
		LabelId: label.EncodeIdToString(model.Id),
	}, nil
}

func (l *LabelServer) ListLabelsForAsset(ctx context.Context, request *label_service.ListLabelsForAssetRequest) (*label_service.ListLabelsForAssetResponse, error) {
	assetId, err := asset.ParseId(request.GetAssetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Asset ID.")
	}

	ids, err := l.LabelManager.ListLabelsForAsset(assetId)
	if err != nil {
		l.Logger.Error("Failed retrieving labels for asset.",
			zap.String("assetId", request.GetAssetId()),
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Internal error.")
	}

	encodedIds := make([]string, len(ids))
	for i, item := range ids {
		encodedIds[i] = label.EncodeIdToString(item)
	}

	return &label_service.ListLabelsForAssetResponse{LabelIds: encodedIds}, nil
}

func (l *LabelServer) GetLabelUrl(ctx context.Context, request *label_service.GetLabelUrlRequest) (*label_service.GetLabelUrlResponse, error) {
	labelId, err := label.ParseId(request.LabelId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid label ID.")
	}
	url, err := l.LabelManager.GetRenderedImageURL(labelId)
	if err != nil {
		if errors.Is(err, label_manager.ErrLabelNotFound) {
			return nil, status.Errorf(codes.NotFound, "Label does not exist.")
		}

		l.Logger.Error("Failed retrieving label URL.", zap.String("labelId", request.LabelId), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Internal error.")
	}

	return &label_service.GetLabelUrlResponse{Url: url}, nil
}
