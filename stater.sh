#!/bin/bash

# root directory (change if needed)
ROOT="go-kit"

# create folders
mkdir -p $ROOT/{core,config,models,infrastructure/aws/dynamo,infrastructure/aws/s3,infrastructure/aws/sqs,infrastructure/aws/cognito,modules/{limitclient,searchclient,authclient,notificationclient},platform/{middleware,logger,cache}}

# create files
touch $ROOT/kit.go

# core
touch $ROOT/core/{httpclient.go,auth.go,errors.go,response.go}

# config
touch $ROOT/config/{config.go,defaults.go}

# models
touch $ROOT/models/{user.go,property.go,subscription.go,common.go}

# infrastructure aws
touch $ROOT/infrastructure/aws/aws.go
touch $ROOT/infrastructure/aws/dynamo/dynamo.go
touch $ROOT/infrastructure/aws/s3/{s3.go,presigned.go}
touch $ROOT/infrastructure/aws/sqs/sqs.go
touch $ROOT/infrastructure/aws/cognito/{cognito.go,context.go}

# modules (empty dirs already created)

# platform
touch $ROOT/platform/middleware/{auth.go,ratelimit.go}
touch $ROOT/platform/logger/logger.go
touch $ROOT/platform/cache/cache.go

echo "Project structure created successfully!"