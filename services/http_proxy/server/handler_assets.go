package server

import (
	core_asset "asset-tracker/pkg/core/asset"
	"asset-tracker/proto/asset_service"
	"asset-tracker/proto/label_service"
	"asset-tracker/services/http_proxy/server/proxy_error"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type listAssetsInput struct {
	StartToken string `form:"StartToken"`
	MaxItems   uint64 `form:"MaxItems"`
}

type asset struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type listAssetsOutput struct {
	Assets    []asset `json:"assets"`
	NextToken *string `json:"nextToken,omitempty"`
}

type getAssetOutput struct {
	Asset asset `json:"asset"`
}

type listLabelsForAssetOutput struct {
	LabelIds []string `json:"labelIds"`
}

func (s *ProxyServerImpl) ListAssets(ctx *gin.Context) {
	var input listAssetsInput
	if err := ctx.ShouldBindQuery(&input); err != nil {
		proxy_error.AbortWithErrorResponse(ctx, proxy_error.BadRequest, "Bad request.")
		return
	}

	var startToken *string = nil
	if len(input.StartToken) > 0 {
		startToken = &input.StartToken
	}
	r, err := s.AssetService.ListAssets(ctx, &asset_service.ListAssetsRequest{
		MaxItems:  input.MaxItems,
		NextToken: startToken,
	})
	if err != nil {
		s.Logger.Error("ListAssets operation failed.", zap.Error(err))
		proxy_error.AbortWithErrorResponse(ctx, proxy_error.InternalError, "Internal server error.")
		return
	}

	assets := make([]asset, len(r.Assets))
	for i, a := range r.Assets {
		assets[i] = asset{
			Id:          a.GetId(),
			Name:        a.GetName(),
			Description: a.GetDescription(),
		}
	}
	ctx.JSON(http.StatusOK, listAssetsOutput{
		Assets:    assets,
		NextToken: r.NextToken,
	})
}

func (s *ProxyServerImpl) GetAsset(ctx *gin.Context) {
	rawID := ctx.Param("id")
	if _, err := core_asset.ParseId(rawID); err != nil {
		proxy_error.AbortWithErrorResponse(ctx, proxy_error.BadRequest, "Invalid asset ID.")
		return
	}

	r, err := s.AssetService.GetAsset(ctx, &asset_service.GetAssetRequest{Id: rawID})
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

	a := r.GetAsset()
	ctx.JSON(http.StatusOK, getAssetOutput{
		Asset: asset{
			Id:          a.GetId(),
			Name:        a.GetName(),
			Description: a.GetDescription(),
		},
	})
}

func (s *ProxyServerImpl) ListLabelsForAsset(ctx *gin.Context) {
	rawID := ctx.Param("id")
	if _, err := core_asset.ParseId(rawID); err != nil {
		proxy_error.AbortWithErrorResponse(ctx, proxy_error.BadRequest, "Invalid asset ID.")
		return
	}

	r, err := s.LabelService.ListLabelsForAsset(ctx, &label_service.ListLabelsForAssetRequest{
		AssetId: rawID,
	})
	if err != nil {
		s.Logger.Error("ListLabelsForAsset operation failed.", zap.String("assetId", rawID), zap.Error(err))
		proxy_error.AbortWithErrorResponse(ctx, proxy_error.InternalError, "Internal error.")
	}

	ctx.JSON(http.StatusOK, listLabelsForAssetOutput{
		LabelIds: r.GetLabelIds(),
	})
}
