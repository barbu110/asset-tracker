package asset_manager

import (
	"asset-tracker/pkg/core/asset"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
)

type DynamoDB struct {
	Client    *dynamodb.Client
	TableName string
	Logger    *zap.Logger
}

func assetKey(id *asset.Id) map[string]types.AttributeValue {
	idBytes, _ := id.MarshalBinary()
	return map[string]types.AttributeValue{
		"Id": &types.AttributeValueMemberB{Value: idBytes},
	}
}

func (d *DynamoDB) CreateAsset(a asset.Asset) error {
	item, err := attributevalue.MarshalMap(a)
	if err != nil {
		d.Logger.Error("Failed to serialize asset.", zap.Error(err))
		return fmt.Errorf("item serialization failed: %w", err)
	}

	if _, err := d.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName:           aws.String(d.TableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(Id)"),
	}); err != nil {
		d.Logger.Error("DynamoDB transaction failed.", zap.Error(err))
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}

func (d *DynamoDB) GetAsset(id *asset.Id) (*asset.Asset, error) {
	o, err := d.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       assetKey(id),
		TableName: aws.String(d.TableName),
	})
	if err != nil {
		return nil, fmt.Errorf("GetItem failed: %w", err)
	}
	if o.Item == nil {
		return nil, ErrAssetNotFound
	}

	var a asset.Asset
	if err := attributevalue.UnmarshalMap(o.Item, &a); err != nil {
		return nil, fmt.Errorf("deserialization failed: %w", err)
	}

	return &a, nil
}

func (d *DynamoDB) ListAssets(params *ListAssetsParams) (data PaginatedData[asset.Asset], err error) {
	//TODO implement me
	panic("implement me")
}
