package config

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type AppConfig struct {
	Region    string
	Endpoint  string
	Queue     string
	AccessKey string
	SecretKey string
}

func Load() AppConfig {
	return AppConfig{
		Region:    os.Getenv("AWS_DEFAULT_REGION"),
		Endpoint:  os.Getenv("SQS_ENDPOINT"),
		Queue:     os.Getenv("QUEUE_NAME"),
		AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}
}

func AWSConfig(ctx context.Context, c AppConfig) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx,
		config.WithRegion(c.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.AccessKey, c.SecretKey, "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
				if c.Endpoint != "" {
					return aws.Endpoint{
						URL:           c.Endpoint,
						PartitionID:   "aws",
						SigningRegion: c.Region,
						Source:        aws.EndpointSourceCustom,
					}, nil
				}
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			},
		)),
	)
}
