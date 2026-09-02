package account_repository

import (
	account_model "github.com/evolution-foundation/evolution-go/pkg/account/model"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(account *account_model.Account) error
	GetByEmail(email string) (*account_model.Account, error)
	GetByApiKey(apiKey string) (*account_model.Account, error)
	GetByID(id string) (*account_model.Account, error)
	Update(account *account_model.Account) error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(account *account_model.Account) error {
	return r.db.Create(account).Error
}

func (r *accountRepository) GetByEmail(email string) (*account_model.Account, error) {
	var acc account_model.Account
	err := r.db.Where("email = ?", email).First(&acc).Error
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *accountRepository) GetByApiKey(apiKey string) (*account_model.Account, error) {
	var acc account_model.Account
	err := r.db.Where("api_key = ?", apiKey).First(&acc).Error
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *accountRepository) GetByID(id string) (*account_model.Account, error) {
	var acc account_model.Account
	err := r.db.Where("id = ?", id).First(&acc).Error
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *accountRepository) Update(account *account_model.Account) error {
	return r.db.Save(account).Error
}
