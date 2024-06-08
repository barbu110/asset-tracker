package main

import (
	"asset-tracker/src/proto"
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
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8000))
	if err != nil {
		logger.Error("Could not start TCP listener.", zap.Error(err))
	}

	srv := assetServer{
		Logger:          logger,
		DynamoDBClient:  dynamodbClient,
		AssetsTableName: "asset-manager-assets",
	}
	grpcServer := grpc.NewServer()
	proto.RegisterAssetServer(grpcServer, &srv)
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("Could not start GRPC server.", zap.Error(err))
	}
}
