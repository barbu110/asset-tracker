package cmd

import (
	_ "asset-tracker/pkg/label/renderer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
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

		logger.Info("Hello!",
			zap.String("endpoint", config.AssetServiceEndpoint),
			zap.String("bind", config.Bind))
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
