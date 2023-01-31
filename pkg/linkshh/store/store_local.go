//go:build local

package store

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Link interface {
	StoreUrl(ctx context.Context, hash, uri string) error
	RetrieveUrlForHash(ctx context.Context, hash string) (string, error)
}

func NewStore() (Link, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		// CHANGE THIS TO us-east-1 TO USE AWS proper
		config.WithRegion("localhost"),
		// Comment the below out if not using localhost
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://dynamo:8000", SigningRegion: "localhost"}, nil
			})),
	)
	if err != nil {
		return linkStorage{}, err
	}

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.Credentials = credentials.NewStaticCredentialsProvider("b59xng", "b2sc6o", "")
	})

	return linkStorage{d: svc}, nil
}
