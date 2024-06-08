package main

import (
	"asset-tracker/src/core/asset"
	"asset-tracker/src/proto"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type assetServer struct {
	Logger          *zap.Logger
	DynamoDBClient  *dynamodb.Client
	AssetsTableName string
}

func (server *assetServer) CreateAsset(ctx context.Context, request *proto.CreateAssetRequest) (*proto.CreateAssetResponse, error) {
	// TODO: Include the custom properties there.
	a := asset.New(request.GetName(), request.GetDescription())

	item, err := attributevalue.MarshalMap(a)
	if err != nil {
		server.Logger.Error("Failed to serialize asset.")
		return nil, status.Errorf(codes.Internal, "Internal service error.")
	}

	_, err = server.DynamoDBClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(server.AssetsTableName),
		Item:      item,
	})
	if err != nil {
		assetIdBytes, _ := a.Id.MarshalBinary()
		server.Logger.Error(
			"Failed to write asset to datastore.",
			zap.ByteString("assetId", assetIdBytes),
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal, "Internal service error.")
	}

	return &proto.CreateAssetResponse{
		Asset: &proto.AssetObject{
			Id:          asset.EncodeIdToString(a.Id),
			Name:        a.Name,
			Description: a.Description,
			Attributes:  []*proto.AssetAttribute{},
		},
	}, nil
}
