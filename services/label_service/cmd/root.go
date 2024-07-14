package cmd

import (
	"asset-tracker/pkg/label/renderer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"os"
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

		r := renderer.RasterRenderer{}
		buf, err := r.Render(&renderer.RenderLabelParams{
			FirstLine:   "Ethernet Cable",
			SecondLine:  "Patch, 3m",
			BarcodeData: []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'},
		})
		if err != nil {
			logger.Fatal("Failed to render label.", zap.Error(err))
		}

		f, err := os.CreateTemp("", "label-*.png")
		if err != nil {
			logger.Fatal("Failed to create tempoerary file.", zap.Error(err))
		}
		defer f.Close()

		_, err = f.Write(buf)
		if err != nil {
			logger.Fatal("Failed writing the label to filesystem.", zap.Error(err))
		}

		logger.Info("Label written to filesystem.", zap.String("path", f.Name()))
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
