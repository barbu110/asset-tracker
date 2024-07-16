package cmd

import (
	"asset-tracker/proto/asset_service"
	"asset-tracker/proto/label_service"
	"asset-tracker/services/http_proxy/server"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var logger *zap.Logger

type Config struct {
	Bind string `mapstructure:"BIND"`

	AssetServiceEndpoint string `mapstructure:"ASSET_SERVICE_ENDPOINT"`
	LabelServiceEndpoint string `mapstructure:"LABEL_SERVICE_ENDPOINT"`
}

var rootCmd = &cobra.Command{
	Use:   "http_proxy",
	Short: "HTTP proxy to expose the backend services in a way convenient to browser frontends.",
	Run: func(_ *cobra.Command, _ []string) {
		c, err := loadConfig()
		if err != nil {
			logger.Fatal("Failed to load configuration.", zap.Error(err))
		}

		logger.Info("Configuration loaded.",
			zap.String("assetServiceEndpoint", c.AssetServiceEndpoint),
			zap.String("labelServiceEndpoint", c.LabelServiceEndpoint))

		httpRouter := gin.Default()

		labelSvcConn, err := grpc.NewClient(c.LabelServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal("Connection to Label Service failed.", zap.Error(err))
		}
		defer func() { _ = labelSvcConn.Close() }()
		labelSvcClient := label_service.NewLabelClient(labelSvcConn)

		assetSvcConn, err := grpc.NewClient(c.AssetServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal("Connection to Asset Service failed.", zap.Error(err))
		}
		defer func() { _ = assetSvcConn.Close() }()
		assetSvcClient := asset_service.NewAssetClient(assetSvcConn)

		proxyServer := server.ProxyServerImpl{
			Logger:       logger.Named("ProxyServer"),
			LabelService: labelSvcClient,
			AssetService: assetSvcClient,
		}
		server.SetupHTTPRouter(&proxyServer, httpRouter)

		srv := &http.Server{
			Addr:    c.Bind,
			Handler: httpRouter,
		}

		go func() {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Fatal("Failed to initialize server.", zap.Error(err))
			}
		}()

		logger.Info("HTTP proxy started.", zap.String("bind", srv.Addr))

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Shutdown server
		logger.Info("Server shutting down.")
		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatal("Server forced to shutdown.", zap.Error(err))
		}
	},
}

func Execute() {
	l, err := zap.NewDevelopment()
	if err != nil {
		log.Panicf("Could not initialize logger: %v", err)
	}
	defer l.Sync()
	logger = l

	rootCmd.Flags().String("bind", "0.0.0.0:8080", "IP & port for the TCP listener.")
	_ = viper.BindPFlag("bind", rootCmd.Flags().Lookup("bind"))

	rootCmd.Flags().String("asset_service_endpoint", "", "Endpoint of the Asset Service.")
	_ = viper.BindPFlag("asset_service_endpoint", rootCmd.Flags().Lookup("asset_service_endpoint"))

	rootCmd.Flags().String("label_service_endpoint", "", "Endpoint of the Label Service.")
	_ = viper.BindPFlag("label_service_endpoint", rootCmd.Flags().Lookup("label_service_endpoint"))

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("Command execution failed.", zap.Error(err))
	}
}

func loadConfig() (Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("alm")
	viper.AllowEmptyEnv(false)
	viper.MustBindEnv(
		"asset_service_endpoint",
		"label_service_endpoint",
	)

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return Config{}, err
	}
	return c, nil
}
