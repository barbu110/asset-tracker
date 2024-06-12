package asset_manager

import (
	"asset-tracker/pkg/core/asset"
	"asset-tracker/pkg/pagination"
	"asset-tracker/pkg/pagination/next_token"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/zap"
)

type DynamoDB struct {
	Client                    *dynamodb.Client
	TableName                 string
	NextTokenEncryptionEngine next_token.EncryptionEngine
	Logger                    *zap.Logger
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
		return fmt.Errorf("item serialization failed: %w", err)
	}

	if _, err := d.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName:           aws.String(d.TableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(Id)"),
	}); err != nil {
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

func (d *DynamoDB) ListAssets(params *ListAssetsParams) (data pagination.PaginatedData[asset.Asset], err error) {
	startKey, err := d.decodeStartKey(params.NextToken, params.HasNextToken)
	if err != nil {
		d.Logger.Debug(
			"Could not parse the received NextToken.",
			zap.String("NextToken", params.NextToken),
			zap.Error(err),
		)
		return pagination.NewEmpty[asset.Asset](), ErrInvalidNextToken
	}

	i := dynamodb.ScanInput{
		TableName:         aws.String(d.TableName),
		Limit:             aws.Int32(int32(params.GetMaxItems())),
		Select:            "ALL_ATTRIBUTES",
		ExclusiveStartKey: startKey,
	}
	output, err := d.Client.Scan(context.TODO(), &i)
	if err != nil {
		return pagination.NewEmpty[asset.Asset](), fmt.Errorf("could not scan datastore: %w", err)
	}

	assets := make([]asset.Asset, len(output.Items))
	for i, item := range output.Items {
		if err := attributevalue.UnmarshalMap(item, &assets[i]); err != nil {
			return pagination.NewEmpty[asset.Asset](), errors.New("deserialization of item failed")
		}
	}

	lastKey, isPresent, err := d.encodeLastKey(output.LastEvaluatedKey)
	if err != nil {
		return pagination.NewEmpty[asset.Asset](), errors.New("nextToken encoding failed")
	}

	return pagination.PaginatedData[asset.Asset]{
		Items:        assets,
		NextToken:    lastKey,
		HasNextToken: isPresent,
	}, nil
}

func (d *DynamoDB) decodeStartKey(raw string, isPresent bool) (map[string]types.AttributeValue, error) {
	if !isPresent {
		return nil, nil
	}

	j, err := d.NextTokenEncryptionEngine.DecryptFromString(raw)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt token: %w", err)
	}

	km := make(map[string]interface{})
	if err := msgpack.Unmarshal(j.Bytes(), &km); err != nil {
		d.Logger.Debug("Could not unmarshal decrypted token.", zap.ByteString("decryptedToken", j.Bytes()))
		return nil, fmt.Errorf("could not unmarshal decrypted token: %w", err)
	}

	m, err := attributevalue.MarshalMap(km)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal AttributeValue map: %w", err)
	}

	return m, nil
}

func (d *DynamoDB) encodeLastKey(k map[string]types.AttributeValue) (encrypted string, isPresent bool, err error) {
	isPresent = true

	if k == nil {
		isPresent = false
		return
	}

	var km map[string]interface{}
	if e := attributevalue.UnmarshalMap(k, &km); e != nil {
		err = fmt.Errorf("unmarshal map failed: %w", e)
		return
	}
	fmt.Printf("got km = %v\n", km)
	m, err := msgpack.Marshal(km)
	if err != nil {
		err = fmt.Errorf("could not marshal token map: %w", err)
		return
	}
	fmt.Printf("got m = %v\n", string(m))

	encrypted, err = d.NextTokenEncryptionEngine.EncryptToString(d.NextTokenEncryptionEngine.NewToken(m))
	if err != nil {
		err = fmt.Errorf("could not encrypt token: %w", err)
		return
	}

	return
}
