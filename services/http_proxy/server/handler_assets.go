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
	Id          string          `json:"id"`
	AssetKind   core_asset.Kind `json:"assetKind"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	ContainerId string          `json:"containerId,omitempty"`
	CreatedAt   string          `json:"createdAt"`
	ModifiedAt  string          `json:"modifiedAt"`
}

type listAssetsOutput struct {
	Assets    []asset `json:"assets"`
	NextToken *string `json:"nextToken,omitempty"`
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
			ContainerId: a.GetContainerId(),
			Name:        a.GetName(),
			AssetKind:   core_asset.Kind(a.GetAssetKind()),
			Description: a.GetDescription(),
		}
	}
	ctx.JSON(http.StatusOK, listAssetsOutput{
		Assets:    assets,
		NextToken: r.NextToken,
	})
}

func (s *ProxyServerImpl) ListAssetsInContainer(ctx *gin.Context) {
	containerID := ctx.Param("containerId")
	if _, err := core_asset.ParseId(containerID); err != nil {
		proxy_error.AbortWithErrorResponse(ctx, proxy_error.BadRequest, "Invalid container ID.")
		return
	}

	var input listAssetsInput
	if err := ctx.ShouldBindQuery(&input); err != nil {
		proxy_error.AbortWithErrorResponse(ctx, proxy_error.BadRequest, "Bad request.")
		return
	}

	var startToken *string = nil
	if len(input.StartToken) > 0 {
		startToken = &input.StartToken
	}
	r, err := s.AssetService.ListAssetsInContainer(ctx, &asset_service.ListAssetsInContainerRequest{
		ContainerId: containerID,
		MaxItems:    input.MaxItems,
		NextToken:   startToken,
	})
	if err != nil {
		s.Logger.Error("ListAssetsInContainer operation failed.", zap.Error(err))
		proxy_error.AbortWithErrorResponse(ctx, proxy_error.InternalError, "Internal server error.")
		return
	}

	assets := make([]asset, len(r.Assets))
	for i, a := range r.Assets {
		assets[i] = asset{
			Id:          a.GetId(),
			ContainerId: a.GetContainerId(),
			Name:        a.GetName(),
			AssetKind:   core_asset.Kind(a.GetAssetKind()),
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

	type outputModel struct {
		Asset asset `json:"asset"`
	}

	ctx.JSON(http.StatusOK, outputModel{
		Asset: asset{
			Id:          a.GetId(),
			ContainerId: a.GetContainerId(),
			AssetKind:   core_asset.Kind(a.GetAssetKind()), // TODO Implement this functionality in the gRPC library.
			Name:        a.GetName(),
			Description: a.GetDescription(),
			CreatedAt:   "",
			ModifiedAt:  "",
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

	type outputModel struct {
		LabelIds []string `json:"labelIds"`
	}

	labelIds := make([]string, 0)
	if r.GetLabelIds() != nil {
		labelIds = r.GetLabelIds()
	}

	ctx.JSON(http.StatusOK, outputModel{
		LabelIds: labelIds,
	})
}
