package account_service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	account_model "github.com/evolution-foundation/evolution-go/pkg/account/model"
	account_repository "github.com/evolution-foundation/evolution-go/pkg/account/repository"
	logger_wrapper "github.com/evolution-foundation/evolution-go/pkg/logger"
	"github.com/evolution-foundation/evolution-go/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type VerifyRequest struct {
	Email  string `json:"email"`
	ApiKey string `json:"apiKey"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	ApiKey   string `json:"apiKey,omitempty"`
}

type AccountResponse struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	ApiKey     string `json:"apiKey"`
	IsVerified bool   `json:"isVerified"`
	CreatedAt  string `json:"createdAt"`
}

type AccountService interface {
	Register(req *RegisterRequest) (*AccountResponse, error)
	Verify(req *VerifyRequest) (*AccountResponse, error)
	Login(req *LoginRequest) (*AccountResponse, error)
	GetByApiKey(apiKey string) (*account_model.Account, error)
}

type accountService struct {
	accountRepo   account_repository.AccountRepository
	loggerWrapper *logger_wrapper.LoggerManager
}

func NewAccountService(accountRepo account_repository.AccountRepository, loggerWrapper *logger_wrapper.LoggerManager) AccountService {
	return &accountService{
		accountRepo:   accountRepo,
		loggerWrapper: loggerWrapper,
	}
}

func generateSecureApiKey() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return "evo_live_" + hex.EncodeToString(b)
}

func (s *accountService) Register(req *RegisterRequest) (*AccountResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" {
		return nil, errors.New("email is required")
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("name is required")
	}
	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters long")
	}

	existing, _ := s.accountRepo.GetByEmail(email)
	if existing != nil {
		return nil, errors.New("an account with this email already exists")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	apiKey := generateSecureApiKey()

	acc := &account_model.Account{
		Email:        email,
		Name:         name,
		PasswordHash: string(hashedBytes),
		ApiKey:       apiKey,
		IsActive:     true,
		IsVerified:   false,
	}

	err = s.accountRepo.Create(acc)
	if err != nil {
		return nil, fmt.Errorf("failed to create user account: %w", err)
	}

	// Trigger async email notification with API Key via Office 365 SMTP
	go func(targetEmail, targetName, targetKey string) {
		errEmail := utils.SendApiKeyEmail(targetEmail, targetName, targetKey)
		if errEmail != nil {
			s.loggerWrapper.GetLogger("system").LogError("Failed to send API Key email to %s: %v", targetEmail, errEmail)
		} else {
			s.loggerWrapper.GetLogger("system").LogInfo("Successfully sent API Key email to %s", targetEmail)
		}
	}(email, name, apiKey)

	return &AccountResponse{
		ID:         acc.Id,
		Email:      acc.Email,
		Name:       acc.Name,
		ApiKey:     acc.ApiKey,
		IsVerified: acc.IsVerified,
		CreatedAt:  acc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *accountService) Verify(req *VerifyRequest) (*AccountResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" {
		return nil, errors.New("email is required")
	}
	apiKey := strings.TrimSpace(req.ApiKey)
	if apiKey == "" {
		return nil, errors.New("API key is required")
	}

	acc, err := s.accountRepo.GetByEmail(email)
	if err != nil || acc == nil {
		return nil, errors.New("account not found")
	}

	if acc.ApiKey != apiKey {
		return nil, errors.New("invalid API key provided")
	}

	acc.IsVerified = true
	err = s.accountRepo.Update(acc)
	if err != nil {
		return nil, fmt.Errorf("failed to update account verification status: %w", err)
	}

	return &AccountResponse{
		ID:         acc.Id,
		Email:      acc.Email,
		Name:       acc.Name,
		ApiKey:     acc.ApiKey,
		IsVerified: acc.IsVerified,
		CreatedAt:  acc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *accountService) Login(req *LoginRequest) (*AccountResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" {
		return nil, errors.New("email is required")
	}
	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	acc, err := s.accountRepo.GetByEmail(email)
	if err != nil || acc == nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(acc.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !acc.IsVerified {
		return nil, errors.New("please verify your account first using the API key sent to your email")
	}

	if req.ApiKey != "" && strings.TrimSpace(req.ApiKey) != acc.ApiKey {
		return nil, errors.New("invalid API Key provided")
	}

	return &AccountResponse{
		ID:         acc.Id,
		Email:      acc.Email,
		Name:       acc.Name,
		ApiKey:     acc.ApiKey,
		IsVerified: acc.IsVerified,
		CreatedAt:  acc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *accountService) GetByApiKey(apiKey string) (*account_model.Account, error) {
	return s.accountRepo.GetByApiKey(apiKey)
}
