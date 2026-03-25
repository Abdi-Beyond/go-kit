package dependency

import (
	"github.com/Abdi-Beyond/go-kit/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
)

type AppDependencies struct {
	DBClient *dynamodb.Client
	Config   *config.Config
	Logger   *zap.Logger
	//Repo      *repository.Repositories
}
