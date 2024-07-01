package server

import (
	"asset-tracker/proto/label_service"
	"context"
	"go.uber.org/zap"
)

type labelServer struct {
	Logger *zap.Logger
}

func (l *labelServer) RenderLabel(ctx context.Context, request *label_service.RenderLabelRequest) (*label_service.RenderLabelResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (l *labelServer) ListLabelsForAsset(ctx context.Context, request *label_service.ListLabelsForAssetRequest) (*label_service.ListLabelsForAssetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (l *labelServer) GetLabelUrl(ctx context.Context, request *label_service.GetLabelUrlRequest) (*label_service.GetLabelUrlResponse, error) {
	//TODO implement me
	panic("implement me")
}
