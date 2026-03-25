package kit

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/Abdi-Beyond/go-kit/config"
	"github.com/Abdi-Beyond/go-kit/core/dependency"
	"github.com/Abdi-Beyond/go-kit/infrastructure/aws"
	awsdynamo "github.com/Abdi-Beyond/go-kit/infrastructure/aws/dynamo"

	limitclient "github.com/Abdi-Beyond/go-kit/modules/limitclient/services"
	limitrepo "github.com/Abdi-Beyond/go-kit/repo"
)

// Kit is the master struct for all modules and AWS clients
type Kit struct {
	Dynamo *dynamodb.Client
	S3     *s3.Client
	SQS    *sqs.Client

	Limit *limitclient.PaymentService
}

// New initializes the Kit with all clients and services
func New(ctx context.Context, cfg *config.Config) (*Kit, error) {
	// Load AWS configuration using your aws.Load function

	awsCfg, err := aws.Load(ctx, aws.Config{
		Region:          cfg.AWSRegion,
		AccessKeyID:     cfg.AWSAccessKeyID,
		SecretAccessKey: cfg.AWSSecretAccessKey,
	})
	if err != nil {
		return nil, err
	}

	// Initialize AWS clients
	dynamoClient := awsdynamo.New(awsCfg)

	DEPS := &dependency.AppDependencies{
		DBClient: dynamoClient,
		Config:   cfg,
	}

	// Initialize LimitClient repository and service
	limitRepo := limitrepo.NewRepositories(DEPS)
	limitService := limitclient.NewPaymentService(limitRepo, DEPS)

	// Assemble Kit
	return &Kit{
		Dynamo: dynamoClient,

		Limit: limitService,
	}, nil
}
