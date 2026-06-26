package controllers

import (
	"net/http"

	"github.com/example/ag-directory-manager/internal/dto"
	"github.com/example/ag-directory-manager/internal/middleware"
	"github.com/example/ag-directory-manager/internal/services"
	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	service *services.DashboardService
}

func NewDashboardController(service *services.DashboardService) *DashboardController {
	return &DashboardController{service: service}
}

func (ctl *DashboardController) Snapshot(c *gin.Context) {
	rawOU, _ := c.Get(middleware.ContextOUScopes)
	ouScopes, _ := rawOU.([]string)

	snapshot, err := ctl.service.Snapshot(c.Request.Context(), ouScopes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.DashboardResponse{Snapshot: snapshot})
}
