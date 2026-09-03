package campaign_repository

import (
	"fmt"
	"strings"
	"time"

	campaign_model "github.com/evolution-foundation/evolution-go/pkg/campaign/model"
	"gorm.io/gorm"
)

type DashboardStats struct {
	TotalUsers        int64   `json:"totalUsers"`
	TotalMessagesSent int64   `json:"totalMessagesSent"`
	TotalMessagesFailed int64 `json:"totalMessagesFailed"`
	TotalReplies      int64   `json:"totalReplies"`
	TotalReactions    int64   `json:"totalReactions"`
	TotalResponded    int64   `json:"totalResponded"`
	ResponseRate      string  `json:"responseRate"`
}

type CampaignRepository interface {
	CreateCampaign(campaign *campaign_model.Campaign) error
	UpdateCampaign(campaign *campaign_model.Campaign) error
	CreateCustomer(customer *campaign_model.CampaignCustomer) error
	GetDashboardStats(instanceID string) (*DashboardStats, error)
	GetCustomers(instanceID string, page, limit int, search, status string) ([]campaign_model.CampaignCustomer, int64, error)
	RecordCustomerReply(instanceID, customerNumber, replyText string) error
	RecordCustomerReaction(instanceID, messageID, reactionEmoji string) error
	GetCampaigns(instanceID string) ([]campaign_model.Campaign, error)
}

type campaignRepository struct {
	db *gorm.DB
}

func NewCampaignRepository(db *gorm.DB) CampaignRepository {
	return &campaignRepository{db: db}
}

func (r *campaignRepository) CreateCampaign(campaign *campaign_model.Campaign) error {
	return r.db.Create(campaign).Error
}

func (r *campaignRepository) UpdateCampaign(campaign *campaign_model.Campaign) error {
	return r.db.Save(campaign).Error
}

func (r *campaignRepository) CreateCustomer(customer *campaign_model.CampaignCustomer) error {
	return r.db.Create(customer).Error
}

func (r *campaignRepository) GetDashboardStats(instanceID string) (*DashboardStats, error) {
	stats := &DashboardStats{}

	query := r.db.Model(&campaign_model.CampaignCustomer{})
	if instanceID != "" {
		query = query.Where("instance_id = ?", instanceID)
	}

	query.Count(&stats.TotalUsers)

	r.db.Model(&campaign_model.CampaignCustomer{}).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if instanceID != "" {
				return db.Where("instance_id = ?", instanceID)
			}
			return db
		}).
		Where("message_status = ?", "sent").
		Count(&stats.TotalMessagesSent)

	r.db.Model(&campaign_model.CampaignCustomer{}).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if instanceID != "" {
				return db.Where("instance_id = ?", instanceID)
			}
			return db
		}).
		Where("message_status = ?", "failed").
		Count(&stats.TotalMessagesFailed)

	r.db.Model(&campaign_model.CampaignCustomer{}).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if instanceID != "" {
				return db.Where("instance_id = ?", instanceID)
			}
			return db
		}).
		Where("has_replied = ?", true).
		Count(&stats.TotalReplies)

	r.db.Model(&campaign_model.CampaignCustomer{}).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if instanceID != "" {
				return db.Where("instance_id = ?", instanceID)
			}
			return db
		}).
		Where("has_reacted = ?", true).
		Count(&stats.TotalReactions)

	r.db.Model(&campaign_model.CampaignCustomer{}).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if instanceID != "" {
				return db.Where("instance_id = ?", instanceID)
			}
			return db
		}).
		Where("has_replied = ? OR has_reacted = ?", true, true).
		Count(&stats.TotalResponded)

	if stats.TotalMessagesSent > 0 {
		rate := (float64(stats.TotalResponded) / float64(stats.TotalMessagesSent)) * 100
		stats.ResponseRate = fmt.Sprintf("%.1f%%", rate)
	} else {
		stats.ResponseRate = "0.0%"
	}

	return stats, nil
}

func (r *campaignRepository) GetCustomers(instanceID string, page, limit int, search, status string) ([]campaign_model.CampaignCustomer, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 500 {
		limit = 50
	}
	offset := (page - 1) * limit

	var customers []campaign_model.CampaignCustomer
	var total int64

	query := r.db.Model(&campaign_model.CampaignCustomer{})
	if instanceID != "" {
		query = query.Where("instance_id = ?", instanceID)
	}

	search = strings.TrimSpace(search)
	if search != "" {
		searchLike := "%" + search + "%"
		query = query.Where("name ILIKE ? OR number ILIKE ?", searchLike, searchLike)
	}

	status = strings.ToLower(strings.TrimSpace(status))
	if status == "sent" || status == "failed" {
		query = query.Where("message_status = ?", status)
	} else if status == "replied" {
		query = query.Where("has_replied = ?", true)
	} else if status == "reacted" {
		query = query.Where("has_reacted = ?", true)
	} else if status == "responded" {
		query = query.Where("has_replied = ? OR has_reacted = ?", true, true)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("sent_at DESC").Offset(offset).Limit(limit).Find(&customers).Error
	if err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (r *campaignRepository) RecordCustomerReply(instanceID, customerNumber, replyText string) error {
	cleanNum := cleanPhoneNumber(customerNumber)
	now := time.Now()

	// Find the latest message sent to this customer
	var cust campaign_model.CampaignCustomer
	q := r.db.Where("number LIKE ?", "%"+cleanNum)
	if instanceID != "" {
		q = q.Where("instance_id = ?", instanceID)
	}
	err := q.Order("sent_at DESC").First(&cust).Error
	if err != nil {
		return err
	}

	cust.HasReplied = true
	cust.ReplyText = replyText
	cust.RespondedAt = &now

	return r.db.Save(&cust).Error
}

func (r *campaignRepository) RecordCustomerReaction(instanceID, messageID, reactionEmoji string) error {
	now := time.Now()

	var cust campaign_model.CampaignCustomer
	q := r.db.Where("message_id = ?", messageID)
	if instanceID != "" {
		q = q.Where("instance_id = ?", instanceID)
	}
	err := q.First(&cust).Error
	if err != nil {
		return err
	}

	cust.HasReacted = true
	cust.ReactionEmoji = reactionEmoji
	cust.RespondedAt = &now

	return r.db.Save(&cust).Error
}

func (r *campaignRepository) GetCampaigns(instanceID string) ([]campaign_model.Campaign, error) {
	var campaigns []campaign_model.Campaign
	q := r.db.Model(&campaign_model.Campaign{})
	if instanceID != "" {
		q = q.Where("instance_id = ?", instanceID)
	}
	err := q.Order("created_at DESC").Find(&campaigns).Error
	return campaigns, err
}

func cleanPhoneNumber(num string) string {
	var b strings.Builder
	for _, ch := range num {
		if ch >= '0' && ch <= '9' {
			b.WriteRune(ch)
		}
	}
	res := b.String()
	if len(res) > 10 {
		return res[len(res)-10:]
	}
	return res
}
