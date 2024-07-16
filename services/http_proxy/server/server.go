package server

import (
	"asset-tracker/proto/asset_service"
	"asset-tracker/proto/label_service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProxyServer interface {
	ListAssets(c *gin.Context)
	GetAsset(c *gin.Context)
}

type ProxyServerImpl struct {
	Logger       *zap.Logger
	LabelService label_service.LabelClient
	AssetService asset_service.AssetClient
}

func SetupHTTPRouter(s ProxyServer, engine *gin.Engine) {
	assets := engine.Group("/assets")
	assets.GET("/", s.ListAssets)
	assets.GET("/:id", s.GetAsset)
}
