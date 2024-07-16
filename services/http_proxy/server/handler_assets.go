package server

import (
	"asset-tracker/proto/asset_service"
	"asset-tracker/services/http_proxy/server/proxy_error"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

func (s *ProxyServerImpl) ListAssets(c *gin.Context) {
	var input listAssetsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		proxy_error.AbortWithErrorResponse(c, proxy_error.BadRequest, "Bad request.")
		return
	}

	var startToken *string = nil
	if len(input.StartToken) > 0 {
		startToken = &input.StartToken
	}
	r, err := s.AssetService.ListAssets(context.TODO(), &asset_service.ListAssetsRequest{
		MaxItems:  input.MaxItems,
		NextToken: startToken,
	})
	if err != nil {
		s.Logger.Error("ListAssets operation failed.", zap.Error(err))
		proxy_error.AbortWithErrorResponse(c, proxy_error.InternalError, "Internal server error.")
	}

	assets := make([]asset, len(r.Assets))
	for i, a := range r.Assets {
		assets[i] = asset{
			Id:          a.GetId(),
			Name:        a.GetName(),
			Description: a.GetDescription(),
		}
	}
	c.JSON(http.StatusOK, listAssetsOutput{
		Assets:    assets,
		NextToken: r.NextToken,
	})
}
