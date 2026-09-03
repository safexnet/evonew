package account_handler

import (
	"net/http"

	account_model "github.com/evolution-foundation/evolution-go/pkg/account/model"
	account_service "github.com/evolution-foundation/evolution-go/pkg/account/service"
	"github.com/gin-gonic/gin"
)

type AccountHandler interface {
	Register(ctx *gin.Context)
	Verify(ctx *gin.Context)
	Login(ctx *gin.Context)
	Me(ctx *gin.Context)
}

type accountHandler struct {
	accountService account_service.AccountService
}

func NewAccountHandler(accountService account_service.AccountService) AccountHandler {
	return &accountHandler{accountService: accountService}
}

// Register handles POST /account/register
// Flow step 1: User submits name, email, and password.
// System hashes password, generates unique API Key (evo_live_...), saves account to DB (isVerified=false),
// sends API key via email, and returns 201 Created with account details.
// @Summary Register a new account
// @Description Register a new user account, save credentials to DB, generate an API key, and send it via email
// @Tags Account
// @Accept json
// @Produce json
// @Param request body account_service.RegisterRequest true "Registration credentials"
// @Success 201 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /account/register [post]
func (h *accountHandler) Register(ctx *gin.Context) {
	var req account_service.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc, err := h.accountService.Register(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "An API key has been sent if this email is valid. Please verify your account using that key.",
		"data":    acc,
	})
}

// Verify handles POST /account/verify
// Flow step 2: User enters registered email and API Key (received via email or creation response).
// System validates key, sets isVerified=true in DB, enabling the user to log in.
// @Summary Verify an account using email and API key
// @Description Verify a new user account using their email and the API key sent via email
// @Tags Account
// @Accept json
// @Produce json
// @Param request body account_service.VerifyRequest true "Verification payload"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /account/verify [post]
func (h *accountHandler) Verify(ctx *gin.Context) {
	var req account_service.VerifyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc, err := h.accountService.Verify(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Account verified successfully! You can now log in and access services.",
		"data":    acc,
	})
}

// Login handles POST /account/login
// Flow step 3: User submits registered email and password.
// System checks password hash and account verification status (isVerified must be true).
// Returns 200 OK with user details including account API Key, which frontend uses for session auth.
// @Summary Log in to an account
// @Description Authenticate user credentials and return API key
// @Tags Account
// @Accept json
// @Produce json
// @Param request body account_service.LoginRequest true "Login credentials"
// @Success 200 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /account/login [post]
func (h *accountHandler) Login(ctx *gin.Context) {
	var req account_service.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc, err := h.accountService.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    acc,
	})
}

// @Summary Get current account profile
// @Description Get current authenticated user profile
// @Tags Account
// @Produce json
// @Param apikey header string true "API Key"
// @Success 200 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /account/me [get]
func (h *accountHandler) Me(ctx *gin.Context) {
	getAcc, exists := ctx.Get("account")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	acc, ok := getAcc.(*account_model.Account)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid account session"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":         acc.Id,
			"email":      acc.Email,
			"name":       acc.Name,
			"apiKey":     acc.ApiKey,
			"isVerified": acc.IsVerified,
			"createdAt":  acc.CreatedAt,
		},
	})
}
