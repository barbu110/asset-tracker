package main

import (
	"asset-tracker/pkg/asset_manager"
	"asset-tracker/pkg/pagination/next_token"
	"asset-tracker/proto/asset_service"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	dynamodbClient := dynamodb.NewFromConfig(cfg)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	address := fmt.Sprintf("localhost:%d", 8000)
	logger.Info("Starting TCP listener", zap.String("address", address))

	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("Could not start TCP listener.", zap.Error(err))
	}

	srv := assetServer{
		Logger: logger,
		AssetManager: &asset_manager.DynamoDB{
			Client:    dynamodbClient,
			TableName: "asset-manager-assets",
			NextTokenEncryptionEngine: next_token.EncryptionEngine{
				KeySource: &next_token.EnvironmentKeySource{VariableName: "NEXT_TOKEN_KEY"},
			},
			Logger: logger.Named("AssetManager"),
		},
	}
	grpcServer := grpc.NewServer()
	asset_service.RegisterAssetServer(grpcServer, &srv)
	reflection.Register(grpcServer)

	logger.Info("Registered AssetService onto the GRPC server.")

	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("Could not start GRPC server.", zap.Error(err))
	}
}
