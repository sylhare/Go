package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
	"log"
	"net"
)

type LocalStackContainer struct {
	Config aws.Config

	*localstack.LocalStackContainer
}

func Localstack() (*LocalStackContainer, error) {
	ctx := context.Background()

	localstackContainer, err := localstack.Run(
		ctx,
		"localstack/localstack:3.8.1",
		testcontainers.WithEnv(map[string]string{
			"SERVICES": "s3,sqs,sns",
			"DEBUG":    "1",
		}),
	)
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return nil, err
	}

	host, err := localstackContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := localstackContainer.MappedPort(ctx, "4566/tcp")
	if err != nil {
		return nil, err
	}

	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		return nil, err
	}

	awsConfig.BaseEndpoint = aws.String("http://" + net.JoinHostPort(host, port.Port()))

	return &LocalStackContainer{
		awsConfig,
		localstackContainer,
	}, nil

}
