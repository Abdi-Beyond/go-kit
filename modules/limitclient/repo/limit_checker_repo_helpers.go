package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Abdi-Beyond/go-kit/core/dependency"
	"github.com/Abdi-Beyond/go-kit/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
	"github.com/sirupsen/logrus"
)

type LimitCheckerRepo struct {
	deps *dependency.AppDependencies
}

func NewLimitCheckerRepo(deps *dependency.AppDependencies) *LimitCheckerRepo {
	return &LimitCheckerRepo{
		deps: deps,
	}
}


func (r *LimitCheckerRepo) CountUserOwnedFolders(ctx context.Context, userID string) (int64, error) {
	pk := fmt.Sprintf("USER#%s", userID)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.deps.Config.DynamoDBTables.Folders),
		KeyConditionExpression: aws.String("PK = :pk"),
		FilterExpression:       aws.String("#rel = :owner"),
		ExpressionAttributeNames: map[string]string{
			"#rel": "Relation",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":    &types.AttributeValueMemberS{Value: pk},
			":owner": &types.AttributeValueMemberS{Value: "owner"},
		},
		Select: types.SelectCount,
	}

	out, err := r.deps.DBClient.Query(ctx, input)
	if err != nil {
		return 0, err
	}

	return int64(out.Count), nil
}

func (r *LimitCheckerRepo) CountUserProperties(ctx context.Context, userID string) (int64, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.deps.Config.DynamoDBTables.UserPropertyCore),
		IndexName:              aws.String("OwnerCreatedAtIndex"),
		KeyConditionExpression: aws.String("userID = :uid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberS{Value: userID},
		},
		Select: types.SelectCount,
	}

	out, err := r.deps.DBClient.Query(ctx, input)
	if err != nil {
		return 0, err
	}

	return int64(out.Count), nil
}
func (r *LimitCheckerRepo) CountSharesForProperty(
	ctx context.Context,
	propertyID string,
) (int64, error) {

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.deps.Config.DynamoDBTables.UserSharedProperty),
		IndexName:              aws.String("propertySharedUsersIndex"),
		KeyConditionExpression: aws.String("propertyID = :pid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pid": &types.AttributeValueMemberS{Value: propertyID},
		},
		Select: types.SelectCount,
	}

	out, err := r.deps.DBClient.Query(ctx, input)
	if err != nil {
		return 0, err
	}

	return int64(out.Count), nil
}

//    --------------

func (r *LimitCheckerRepo) PlanDefinition(ctx context.Context, planID string, feature string) (*models.PlanFeature, error) {
	// small letter

	logrus.Printf("planID: %s", planID)
	logrus.Printf("feature: %s", feature)

	smallFeature := strings.ToLower(planID)
	pk := fmt.Sprintf("PLAN#%s", smallFeature)

	sk := fmt.Sprintf("FEATURE#%s", feature)

	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.deps.Config.DynamoDBTables.PlanFeatures),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{
				Value: pk,
			},
			"SK": &types.AttributeValueMemberS{
				Value: sk,
			},
		},
	}

	result, err := r.deps.DBClient.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("plan not found")
	}

	var plan models.PlanFeature

	err = attributevalue.UnmarshalMap(result.Item, &plan)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}

func (r *LimitCheckerRepo) SeedPlanFeatures(ctx context.Context, req models.PlanFeature) error {
	now := time.Now().UnixMilli()

	// set audit fields if not provided
	if req.CreatedAt == 0 {
		req.CreatedAt = now
	}
	req.UpdatedAt = now

	// marshal to DynamoDB format
	item, err := attributevalue.MarshalMap(req)
	if err != nil {
		return fmt.Errorf("failed to marshal plan feature: %w", err)
	}

	_, err = r.deps.DBClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.deps.Config.DynamoDBTables.PlanFeatures),
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("failed to put plan feature: %w", err)
	}

	return nil
}
func (r *LimitCheckerRepo) GetByUserID(ctx context.Context, userID string) (*models.SubscriptionList, error) {
	if r.deps.Config.DynamoDBTables.Subscriptions == "" {
		return nil, fmt.Errorf("DynamoDB table name is empty")
	}
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.deps.Config.DynamoDBTables.Subscriptions),
		IndexName:              aws.String("subscriptions-userid-index"),
		KeyConditionExpression: aws.String("user_id = :uid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberS{Value: userID},
		},
		Limit:            aws.Int32(1),    // Only one item
		ScanIndexForward: aws.Bool(false), // Get most recent first
	}

	resp, err := r.deps.DBClient.Query(ctx, input)
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case "ValidationException":
				return nil, fmt.Errorf("query validation failed: %w", err)
			case "ResourceNotFoundException":
				return nil, fmt.Errorf("table or index not found: %w", err)
			}
		}
		return nil, fmt.Errorf("failed to query user subscription: %w", err)
	}

	if len(resp.Items) == 0 {
		return nil, nil // No subscription found for this user
	}

	var sub models.SubscriptionList
	if err := attributevalue.UnmarshalMap(resp.Items[0], &sub); err != nil {
		return nil, fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	return &sub, nil
}
