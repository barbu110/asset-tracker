package server

import (
	"asset-tracker/proto/asset_service"
	"asset-tracker/proto/label_service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProxyServer interface {
	CreateAsset(c *gin.Context)
	ListAssets(c *gin.Context)
	ListAssetsInContainer(c *gin.Context)
	GetAsset(c *gin.Context)
	ListLabelsForAsset(c *gin.Context)
	RenderLabel(c *gin.Context)
	MakeLabelUrl(c *gin.Context)
}

type ProxyServerImpl struct {
	Logger       *zap.Logger
	LabelService label_service.LabelClient
	AssetService asset_service.AssetClient
}

func SetupHTTPRouter(s ProxyServer, engine *gin.Engine) {
	assets := engine.Group("/assets")
	assets.GET("/", s.ListAssets)
	assets.POST("/", s.CreateAsset)
	assets.GET("/:id", s.GetAsset)
	assets.GET("/:id/labels", s.ListLabelsForAsset)

	engine.GET("/assets-in-container/:containerId", s.ListAssetsInContainer)

	labels := engine.Group("/labels")
	labels.POST("/:labelId", s.MakeLabelUrl)
	labels.POST("/render/:assetId", s.RenderLabel)
}
