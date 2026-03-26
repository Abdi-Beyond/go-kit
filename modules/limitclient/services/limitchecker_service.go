package service

import (
	"context"
	"errors"

	"github.com/Abdi-Beyond/go-kit/core/dependency"
	"github.com/Abdi-Beyond/go-kit/models"
	"github.com/Abdi-Beyond/go-kit/modules/limitclient/dto"
	"github.com/Abdi-Beyond/go-kit/repo"
)

type PaymentService struct {
	Repo *repo.Repositories
	Deps *dependency.AppDependencies
}

// Updated constructor to accept the DynamoDB client
func NewPaymentService(paymentRepo *repo.Repositories, deps *dependency.AppDependencies) *PaymentService {
	return &PaymentService{
		Repo: paymentRepo,
		Deps: deps,
	}
}

func (s *PaymentService) FoldersNumberPerUser(ctx context.Context, userid string) (int64, error) {
	return s.Repo.Limitcheck.CountUserOwnedFolders(ctx, userid)
}

func (s *PaymentService) DealNumberonProperty(ctx context.Context, userid, propertyid string) (int64, error) {
	return s.Repo.Limitcheck.CountSharesForProperty(ctx, propertyid)
}

func (s *PaymentService) PropertyCreatedBytheUser(ctx context.Context, userid string) (int64, error) {
	return s.Repo.Limitcheck.CountUserProperties(ctx, userid)
}

func (s *PaymentService) PlanDefinition(ctx context.Context, userid string, feature string) (*models.PlanFeature, error) {
	sub, err := s.Repo.Limitcheck.GetByUserID(ctx, userid)
	if err != nil {
		return nil, err
	}

	plans, err := s.Repo.Limitcheck.PlanDefinition(ctx, sub.PlanID, feature)
	if err != nil {
		return nil, err
	}
	return plans, nil
}

func (s *PaymentService) SeedPlanFeatures(ctx context.Context, req models.PlanFeature) error {
	return s.Repo.Limitcheck.SeedPlanFeatures(ctx, req)
}

// ---  --- ---

func (s *PaymentService) CheckLimit(ctx context.Context, userid string, req dto.LimitCheckRequest) (int64, error) {
	switch req.ServiceType {
	case "folder_creation":
		return s.FoldersNumberPerUser(ctx, userid)
	case "property_shared":
		if req.PropertyID == "" {
			return 0, errors.New("property_id is required for deals")
		}
		return s.DealNumberonProperty(ctx, userid, req.PropertyID)

	default:
		return 0, errors.New("invalid checked_service")
	}

}

func (s *PaymentService) Is_Eligible(ctx context.Context, userid string, feature string) (bool, error) {
	// Get the plan's defined limit for this feature
	planFeature, err := s.PlanDefinition(ctx, userid, feature)
	if err != nil {
		return false, err
	}

	// Get the user's current usage for this feature
	currentUsage, err := s.CheckLimit(ctx, userid, dto.LimitCheckRequest{
		ServiceType: feature,
	})
	if err != nil {
		return false, err
	}

	// Compare current usage against the plan limit
	return currentUsage < *planFeature.Limit, nil
}
