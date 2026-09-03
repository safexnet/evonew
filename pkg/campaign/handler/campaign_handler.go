package campaign_handler

import (
	"net/http"
	"strconv"

	campaign_service "github.com/evolution-foundation/evolution-go/pkg/campaign/service"
	instance_model "github.com/evolution-foundation/evolution-go/pkg/instance/model"
	"github.com/gin-gonic/gin"
)

type CampaignHandler interface {
	GetDashboardStats(ctx *gin.Context)
	GetCustomers(ctx *gin.Context)
	GetCampaigns(ctx *gin.Context)
}

type campaignHandler struct {
	campaignService campaign_service.CampaignService
}

func NewCampaignHandler(campaignService campaign_service.CampaignService) CampaignHandler {
	return &campaignHandler{
		campaignService: campaignService,
	}
}

func getInstanceID(ctx *gin.Context) string {
	instID := ctx.Query("instanceId")
	if instID != "" {
		return instID
	}
	getInst, exists := ctx.Get("instance")
	if exists {
		if inst, ok := getInst.(*instance_model.Instance); ok && inst != nil {
			return inst.Id
		}
	}
	return ""
}

// @Summary Get Campaign Dashboard Statistics
// @Description Get overall analytics metrics (total users, sent, failed, replies, reactions, response rate) for frontend dashboard cards
// @Tags Campaign
// @Produce json
// @Param apikey header string true "API Key"
// @Param instanceId query string false "Instance ID filter (optional)"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /campaign/dashboard [get]
func (h *campaignHandler) GetDashboardStats(ctx *gin.Context) {
	instanceID := getInstanceID(ctx)

	stats, err := h.campaignService.GetDashboardStats(instanceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   stats,
	})
}

// @Summary Get Customer List for Dashboard Table
// @Description Get paginated and searchable list of customers uploaded from Excel with their delivery, reply, and reaction status
// @Tags Campaign
// @Produce json
// @Param apikey header string true "API Key"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 50)"
// @Param search query string false "Search by customer name or phone number"
// @Param status query string false "Filter status: sent, failed, replied, reacted, responded"
// @Param instanceId query string false "Instance ID filter (optional)"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /campaign/customers [get]
func (h *campaignHandler) GetCustomers(ctx *gin.Context) {
	instanceID := getInstanceID(ctx)

	page := 1
	if p, err := strconv.Atoi(ctx.Query("page")); err == nil && p > 0 {
		page = p
	}

	limit := 50
	if l, err := strconv.Atoi(ctx.Query("limit")); err == nil && l > 0 {
		limit = l
	}

	search := ctx.Query("search")
	status := ctx.Query("status")

	result, err := h.campaignService.GetCustomers(instanceID, page, limit, search, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

// @Summary Get All Campaigns
// @Description Get list of bulk message campaigns uploaded via Excel
// @Tags Campaign
// @Produce json
// @Param apikey header string true "API Key"
// @Param instanceId query string false "Instance ID filter (optional)"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /campaign/list [get]
func (h *campaignHandler) GetCampaigns(ctx *gin.Context) {
	instanceID := getInstanceID(ctx)

	campaigns, err := h.campaignService.GetCampaigns(instanceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   campaigns,
	})
}
