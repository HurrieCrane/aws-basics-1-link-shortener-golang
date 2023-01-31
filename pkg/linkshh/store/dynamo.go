package store

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	DynamoDBTableName = "shortened-links"
)

type linkStorage struct {
	d *dynamodb.Client
}

type dynamoModel struct {
	Hash string `json:"link-hash"`
	Link string `json:"link"`
}

func (l linkStorage) StoreUrl(ctx context.Context, hash, uri string) error {
	_, err := l.d.PutItem(ctx, &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"link-hash": &types.AttributeValueMemberS{Value: hash},
			"link":      &types.AttributeValueMemberS{Value: uri},
		},
		TableName: aws.String(DynamoDBTableName),
	})
	if err != nil {
		return err
	}

	return nil
}

func (l linkStorage) RetrieveUrlForHash(ctx context.Context, hash string) (string, error) {
	out, err := l.d.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"link-hash": &types.AttributeValueMemberS{Value: hash},
		},
		TableName: aws.String(DynamoDBTableName),
	})
	if err != nil {
		return "", err
	}

	model := ""
	err = attributevalue.Unmarshal(out.Item["link"], &model)
	if err != nil {
		return "", err
	}

	return model, nil
}
