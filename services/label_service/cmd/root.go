package cmd

import (
	"asset-tracker/proto/label_service"
	"asset-tracker/services/label_service/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

var logger *zap.Logger

type Config struct {
	AssetServiceEndpoint string `mapstructure:"ASSET_SERVICE_ENDPOINT"`
	Bind                 string `mapstructure:"BIND"`
}

var rootCmd = &cobra.Command{
	Use:   "label_service",
	Short: "Service specialized on rendering asset labels.",
	Run: func(_ *cobra.Command, _ []string) {
		config, err := loadConfig()
		if err != nil {
			logger.Fatal("Failed to load configuration.", zap.Error(err))
		}

		srv := server.LabelServer{Logger: logger}
		grpcServer := grpc.NewServer()
		label_service.RegisterLabelServer(grpcServer, &srv)
		reflection.Register(grpcServer)

		logger.Info("Registered LabelService onto the GRPC server.")
		logger.Info("Starting TCP listener", zap.String("address", config.Bind))

		lis, err := net.Listen("tcp", config.Bind)
		if err != nil {
			logger.Error("Could not start TCP listener.", zap.Error(err))
		}

		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("Could not start GRPC server.", zap.Error(err))
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

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("Command execution failed.", zap.Error(err))
	}
}

func loadConfig() (Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("alm")
	viper.AllowEmptyEnv(false)
	viper.MustBindEnv("asset_service_endpoint")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}
	return config, nil
}
