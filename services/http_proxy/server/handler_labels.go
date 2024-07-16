package server

import (
	coreasset "asset-tracker/pkg/core/asset"
	"asset-tracker/pkg/label"
	"asset-tracker/proto/asset_common"
	"asset-tracker/proto/asset_service"
	"asset-tracker/proto/label_service"
	"asset-tracker/services/http_proxy/server/proxy_error"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type labelUrlOutput struct {
	LabelId     string `json:"labelId"`
	RenderedUrl string `json:"renderedUrl"`
}

func (s *ProxyServerImpl) MakeLabelUrl(ctx *gin.Context) {
	rawID := ctx.Param("labelId")
	if _, err := label.ParseId(rawID); err != nil {
		proxy_error.AbortWithErrorResponse(ctx, proxy_error.BadRequest, "Invalid label ID.")
		return
	}

	urlResponse, err := s.LabelService.GetLabelUrl(ctx, &label_service.GetLabelUrlRequest{
		LabelId:     rawID,
		LabelFormat: label_service.LabelFormat_LABEL_FORMAT_PNG,
	})
	if err != nil {
		s.Logger.Error("GetLabelUrl operation failed.",
			zap.String("labelId", rawID),
			zap.Error(err))
		proxy_error.AbortWithErrorResponse(ctx, proxy_error.InternalError, "Internal error.")
	}

	ctx.JSON(http.StatusOK, labelUrlOutput{
		LabelId:     rawID,
		RenderedUrl: urlResponse.GetUrl(),
	})
}

func (s *ProxyServerImpl) RenderLabel(ctx *gin.Context) {
	rawID := ctx.Param("assetId")
	if _, err := coreasset.ParseId(rawID); err != nil {
		proxy_error.AbortWithErrorResponse(ctx, proxy_error.BadRequest, "Invalid asset ID.")
		return
	}

	getAssetResponse, err := s.AssetService.GetAsset(ctx, &asset_service.GetAssetRequest{
		Id: rawID,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				proxy_error.AbortWithErrorResponse(ctx, proxy_error.BadRequest, "Bad request.")
				return
			case codes.NotFound:
				proxy_error.AbortWithErrorResponse(ctx, proxy_error.NotFound, "Asset not found.")
				return
			}
		}

		s.Logger.Error("GetAsset operation failed.", zap.Error(err))
		proxy_error.AbortWithErrorResponse(ctx, proxy_error.InternalError, "Internal error.")
	}

	renderResponse, err := s.LabelService.RenderLabel(ctx, &label_service.RenderLabelRequest{
		Asset: &asset_common.AssetObject{
			Id:          getAssetResponse.GetAsset().GetId(),
			Name:        getAssetResponse.GetAsset().GetName(),
			Description: getAssetResponse.GetAsset().GetDescription(),
		},
	})
	if err != nil {
		s.Logger.Error("RenderLabel operation failed.", zap.Error(err))
		proxy_error.AbortWithErrorResponse(ctx, proxy_error.InternalError, "Internal error.")
		return
	}

	urlResponse, err := s.LabelService.GetLabelUrl(ctx, &label_service.GetLabelUrlRequest{
		LabelId:     renderResponse.GetLabelId(),
		LabelFormat: label_service.LabelFormat_LABEL_FORMAT_PNG,
	})
	if err != nil {
		s.Logger.Error("GetLabelUrl operation failed.",
			zap.String("labelId", renderResponse.GetLabelId()),
			zap.Error(err))
		proxy_error.AbortWithErrorResponse(ctx, proxy_error.InternalError, "Internal error.")
	}

	ctx.JSON(http.StatusOK, labelUrlOutput{
		LabelId:     renderResponse.GetLabelId(),
		RenderedUrl: urlResponse.GetUrl(),
	})
}
