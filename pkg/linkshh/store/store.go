//go:build !local

package store

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Link interface {
	StoreUrl(ctx context.Context, hash, uri string) error
	RetrieveUrlForHash(ctx context.Context, hash string) (string, error)
}

func NewStore() (Link, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return linkStorage{}, err
	}

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	return linkStorage{d: svc}, nil
}
