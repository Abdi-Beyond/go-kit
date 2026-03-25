package kit

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/Abdi-Beyond/go-kit/core/dependency"
	"github.com/Abdi-Beyond/go-kit/infrastructure/aws"
	awsdynamo "github.com/Abdi-Beyond/go-kit/infrastructure/aws/dynamo"
	awss3 "github.com/Abdi-Beyond/go-kit/infrastructure/aws/s3"
	awssqs "github.com/Abdi-Beyond/go-kit/infrastructure/aws/sqs"

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
func New(ctx context.Context, deps *dependency.AppDependencies) (*Kit, error) {
	// Load AWS configuration using your aws.Load function
	awsCfg, err := aws.Load(ctx, aws.Config{
		Region:          deps.Config.AWSRegion,
		AccessKeyID:     deps.Config.AWSAccessKeyID,
		SecretAccessKey: deps.Config.AWSSecretAccessKey,
	})
	if err != nil {
		return nil, err
	}

	// Initialize AWS clients
	dynamoClient := awsdynamo.New(awsCfg)
	s3Client := awss3.New(awsCfg)
	sqsClient := awssqs.New(awsCfg)

	// Initialize LimitClient repository and service
	limitRepo := limitrepo.NewRepositories(deps)
	limitService := limitclient.NewPaymentService(limitRepo, deps)

	// Assemble Kit
	return &Kit{
		Dynamo: dynamoClient,
		S3:     s3Client,
		SQS:    sqsClient,
		Limit:  limitService,
	}, nil
}
