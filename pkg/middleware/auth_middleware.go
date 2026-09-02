package auth_middleware

import (
	"net/http"

	account_service "github.com/evolution-foundation/evolution-go/pkg/account/service"
	"github.com/evolution-foundation/evolution-go/pkg/config"
	instance_service "github.com/evolution-foundation/evolution-go/pkg/instance/service"
	"github.com/gin-gonic/gin"
)

type Middleware interface {
	Auth(ctx *gin.Context)
	AuthAdmin(ctx *gin.Context)
}

type middleware struct {
	config          *config.Config
	instanceService instance_service.InstanceService
	accountService  account_service.AccountService
}

func (m middleware) Auth(ctx *gin.Context) {
	token := ctx.GetHeader("apikey")
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}

	instance, err := m.instanceService.GetInstanceByToken(token)
	if err != nil {
		if token == m.config.GlobalApiKey {
			instances, errAll := m.instanceService.GetAll()
			if errAll == nil && len(instances) > 0 {
				ctx.Set("instance", instances[0])
				ctx.Next()
				return
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no whatsapp instance found. Please create an instance first via /instance/create or in the Manager portal"})
			return
		}

		if m.accountService != nil {
			acc, errAcc := m.accountService.GetByApiKey(token)
			if errAcc == nil && acc != nil {
				ctx.Set("account", acc)
				instances, errAll := m.instanceService.GetAll()
				if errAll == nil && len(instances) > 0 {
					ctx.Set("instance", instances[0])
				}
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}

	ctx.Set("instance", instance)
	ctx.Next()
}

func (m middleware) AuthAdmin(ctx *gin.Context) {
	token := ctx.GetHeader("apikey")
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}

	if token != m.config.GlobalApiKey {
		if m.accountService != nil {
			acc, errAcc := m.accountService.GetByApiKey(token)
			if errAcc == nil && acc != nil {
				ctx.Set("account", acc)
				ctx.Next()
				return
			}
		}
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}

	ctx.Next()
}

func NewMiddleware(config *config.Config, instanceService instance_service.InstanceService, accountService account_service.AccountService) *middleware {
	return &middleware{config: config, instanceService: instanceService, accountService: accountService}
}
