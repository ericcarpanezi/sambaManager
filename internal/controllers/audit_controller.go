package controllers

import (
	"net/http"
	"strconv"

	"github.com/example/ag-directory-manager/internal/audit"
	"github.com/gin-gonic/gin"
)

type AuditController struct {
	service *audit.Service
}

func NewAuditController(service *audit.Service) *AuditController {
	return &AuditController{service: service}
}

func (ctl *AuditController) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	logs, err := ctl.service.List(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": logs})
}
