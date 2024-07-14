package label_manager

import (
	"asset-tracker/pkg/label"
	"asset-tracker/pkg/label/renderer"
	"asset-tracker/pkg/rendered_label_storage"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
	"time"
)

type DynamoDB struct {
	Logger               *zap.Logger
	Client               *dynamodb.Client
	TableName            string
	RenderedLabelStorage rendered_label_storage.RenderedLabelStorage
	LabelRenderer        renderer.LabelRenderer
}

func labelKey(id *label.Id) map[string]types.AttributeValue {
	idBytes, _ := id.MarshalBinary()
	return map[string]types.AttributeValue{
		"Id": &types.AttributeValueMemberB{Value: idBytes},
	}
}

func (d *DynamoDB) CreateLabel(params *CreateLabelParams) (*label.Label, error) {
	labelId := label.RandomId()
	model := label.Label{
		Id:         labelId,
		AssetId:    params.AssetId,
		FirstLine:  params.FirstLine,
		SecondLine: params.SecondLine,
	}
	item, err := attributevalue.MarshalMap(model)
	if err != nil {
		return nil, fmt.Errorf("item serialization failed: %w", err)
	}

	storageResult := make(chan error)
	go func() {
		rendered, err := d.LabelRenderer.Render(&renderer.RenderLabelParams{
			FirstLine:   params.FirstLine,
			SecondLine:  params.SecondLine,
			BarcodeData: params.AssetId,
		})
		if err != nil {
			storageResult <- fmt.Errorf("failed to render label: %w", err)
			return
		}

		if err := d.RenderedLabelStorage.Create(renderedImagePNGKey(labelId), bytes.NewBuffer(rendered)); err != nil {
			storageResult <- fmt.Errorf("failed to store rendered image: %w", err)
			return
		}

		storageResult <- nil
	}()

	ddbResult := make(chan error)
	go func() {
		_, err := d.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName:           aws.String(d.TableName),
			Item:                item,
			ConditionExpression: aws.String("attribute_not_exists(Id)"),
		})
		ddbResult <- err
	}()

	if err := <-ddbResult; err != nil {
		d.Logger.Error("Failed to store label in DynamoDB.", zap.Error(err))
		return nil, fmt.Errorf("failed to label metadata: %w", err)
	}
	if err := <-storageResult; err != nil {
		d.Logger.Error("Failed to store rendered image.", zap.Error(err))
		return nil, fmt.Errorf("failed to store rendered label: %w", err)
	}

	return &model, nil
}

func (d *DynamoDB) GetLabel(id label.Id) (*label.Label, error) {
	o, err := d.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       labelKey(&id),
		TableName: aws.String(d.TableName),
	})
	if err != nil {
		d.Logger.Error("GetItem call failed.", zap.Error(err))
		return nil, fmt.Errorf("GetItem failed: %w", err)
	}
	if o.Item == nil {
		return nil, ErrLabelNotFound
	}

	var l label.Label
	if err := attributevalue.UnmarshalMap(o.Item, &l); err != nil {
		return nil, fmt.Errorf("deserialization failed: %w", err)
	}

	return &l, nil
}

func (d *DynamoDB) GetRenderedImageURL(id label.Id) (string, error) {
	if _, err := d.GetLabel(id); errors.Is(err, ErrLabelNotFound) {
		return "", nil
	}

	url, err := d.RenderedLabelStorage.GetDownloadUrl(renderedImagePNGKey(id), 60*time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed creating download URL: %w", err)
	}

	return url, nil
}
