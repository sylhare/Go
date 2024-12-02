package localstack

import (
	"context"
	sdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
	"log"
	"net"
)

type Localstack struct {
	Config    sdk.Config
	Container *localstack.LocalStackContainer
}

func New() (*Localstack, error) {
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

	awsConfig.BaseEndpoint = sdk.String("http://" + net.JoinHostPort(host, port.Port()))

	return &Localstack{
		awsConfig,
		localstackContainer,
	}, nil
}

func (ls *Localstack) Terminate() {
	if err := testcontainers.TerminateContainer(ls.Container); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
}
