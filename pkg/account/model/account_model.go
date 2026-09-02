package account_model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	Id           string    `json:"id" gorm:"type:uuid;primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-" gorm:"not null"`
	ApiKey       string    `json:"apiKey" gorm:"uniqueIndex;not null"`
	IsActive     bool      `json:"isActive" gorm:"default:true"`
	IsVerified   bool      `json:"isVerified" gorm:"default:false"`
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	if a.Id == "" {
		a.Id = uuid.New().String()
	}
	return
}
