package cmd

import (
	"asset-tracker/pkg/label/renderer"
	"asset-tracker/pkg/label_manager"
	"asset-tracker/pkg/rendered_label_storage"
	"asset-tracker/proto/label_service"
	"asset-tracker/services/label_service/server"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	AssetServiceEndpoint     string `mapstructure:"ASSET_SERVICE_ENDPOINT"`
	Bind                     string `mapstructure:"BIND"`
	LabelsTableName          string `mapstructure:"LABELS_TABLE_NAME"`
	RenderedLabelsBucketName string `mapstructure:"RENDERED_LABELS_BUCKET_NAME"`
}

var rootCmd = &cobra.Command{
	Use:   "label_service",
	Short: "Service specialized on rendering asset labels.",
	Run: func(_ *cobra.Command, _ []string) {
		c, err := loadConfig()
		if err != nil {
			logger.Fatal("Failed to load configuration.", zap.Error(err))
		}

		awsConfig, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			panic(err)
		}

		dynamodbClient := dynamodb.NewFromConfig(awsConfig)
		s3Client := s3.NewFromConfig(awsConfig)
		s3PresignClient := s3.NewPresignClient(s3Client)

		logger.Info("Configuration loaded.",
			zap.String("labelsTableName", c.LabelsTableName),
			zap.String("renderedLabelsBucketName", c.RenderedLabelsBucketName))

		labelStorage := rendered_label_storage.S3{
			Logger:        logger,
			Client:        s3Client,
			PresignClient: s3PresignClient,
			BucketName:    c.RenderedLabelsBucketName,
		}
		labelManager := label_manager.DynamoDB{
			Logger:               logger,
			Client:               dynamodbClient,
			TableName:            c.LabelsTableName,
			RenderedLabelStorage: &labelStorage,
			LabelRenderer:        &renderer.RasterRenderer{},
		}
		srv := server.LabelServer{Logger: logger, LabelManager: &labelManager}
		grpcServer := grpc.NewServer()
		label_service.RegisterLabelServer(grpcServer, &srv)
		reflection.Register(grpcServer)

		logger.Info("Registered LabelService onto the GRPC server.")
		logger.Info("Starting TCP listener", zap.String("address", c.Bind))

		lis, err := net.Listen("tcp", c.Bind)
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

	rootCmd.Flags().String("labels_table_name", "", "Name of DynamoDB table holding labels data.")
	_ = viper.BindPFlag("labels_table_name", rootCmd.Flags().Lookup("labels_table_name"))

	rootCmd.Flags().String("rendered_labels_bucket_name", "", "Name of S3 bucket holding rendered labels.")
	_ = viper.BindPFlag("rendered_labels_bucket_name", rootCmd.Flags().Lookup("rendered_labels_bucket_name"))

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
		"labels_table_name",
		"rendered_labels_bucket_name",
	)

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return Config{}, err
	}
	return c, nil
}
