package campaign_service

import (
	campaign_model "github.com/evolution-foundation/evolution-go/pkg/campaign/model"
	campaign_repository "github.com/evolution-foundation/evolution-go/pkg/campaign/repository"
	logger_wrapper "github.com/evolution-foundation/evolution-go/pkg/logger"
)

type CustomerListResponse struct {
	Total     int64                              `json:"total"`
	Page      int                                `json:"page"`
	Limit     int                                `json:"limit"`
	Customers []campaign_model.CampaignCustomer `json:"customers"`
}

type CampaignService interface {
	GetDashboardStats(instanceID string) (*campaign_repository.DashboardStats, error)
	GetCustomers(instanceID string, page, limit int, search, status string) (*CustomerListResponse, error)
	GetCampaigns(instanceID string) ([]campaign_model.Campaign, error)
	GetRepository() campaign_repository.CampaignRepository
}

type campaignService struct {
	campaignRepo  campaign_repository.CampaignRepository
	loggerWrapper *logger_wrapper.LoggerManager
}

func NewCampaignService(campaignRepo campaign_repository.CampaignRepository, loggerWrapper *logger_wrapper.LoggerManager) CampaignService {
	return &campaignService{
		campaignRepo:  campaignRepo,
		loggerWrapper: loggerWrapper,
	}
}

func (s *campaignService) GetDashboardStats(instanceID string) (*campaign_repository.DashboardStats, error) {
	return s.campaignRepo.GetDashboardStats(instanceID)
}

func (s *campaignService) GetCustomers(instanceID string, page, limit int, search, status string) (*CustomerListResponse, error) {
	customers, total, err := s.campaignRepo.GetCustomers(instanceID, page, limit, search, status)
	if err != nil {
		return nil, err
	}
	return &CustomerListResponse{
		Total:     total,
		Page:      page,
		Limit:     limit,
		Customers: customers,
	}, nil
}

func (s *campaignService) GetCampaigns(instanceID string) ([]campaign_model.Campaign, error) {
	return s.campaignRepo.GetCampaigns(instanceID)
}

func (s *campaignService) GetRepository() campaign_repository.CampaignRepository {
	return s.campaignRepo
}
