package campaign_model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Campaign struct {
	ID         string    `json:"id" gorm:"type:uuid;primaryKey"`
	Name       string    `json:"name"`
	InstanceID string    `json:"instanceId" gorm:"index"`
	Total      int       `json:"total"`
	Sent       int       `json:"sent"`
	Failed     int       `json:"failed"`
	CreatedAt  time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (c *Campaign) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}

type CampaignCustomer struct {
	ID            string     `json:"id" gorm:"type:uuid;primaryKey"`
	CampaignID    string     `json:"campaignId" gorm:"type:uuid;index"`
	CampaignName  string     `json:"campaignName"`
	InstanceID    string     `json:"instanceId" gorm:"index"`
	Name          string     `json:"name"`
	Number        string     `json:"number" gorm:"index"`
	MessageID     string     `json:"messageId" gorm:"index"`
	MessageText   string     `json:"messageText"`
	MessageStatus string     `json:"messageStatus" gorm:"index"` // "sent", "failed"
	ErrorMessage  string     `json:"errorMessage,omitempty"`
	HasReplied    bool       `json:"hasReplied" gorm:"default:false;index"`
	ReplyText     string     `json:"replyText,omitempty"`
	HasReacted    bool       `json:"hasReacted" gorm:"default:false;index"`
	ReactionEmoji string     `json:"reactionEmoji,omitempty"`
	RespondedAt   *time.Time `json:"respondedAt,omitempty"`
	SentAt        time.Time  `json:"sentAt"`
	CreatedAt     time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (c *CampaignCustomer) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}
