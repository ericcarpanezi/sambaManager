package controllers

import (
	"net/http"

	"github.com/example/ag-directory-manager/internal/services"
	"github.com/gin-gonic/gin"
)

type SettingsController struct {
	service *services.SettingsService
}

func NewSettingsController(service *services.SettingsService) *SettingsController {
	return &SettingsController{service: service}
}

func (ctl *SettingsController) GetLDAP(c *gin.Context) {
	settings, err := ctl.service.GetLDAP(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, settings)
}

func (ctl *SettingsController) SaveLDAP(c *gin.Context) {
	var payload services.LDAPSettings
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctl.service.SaveLDAP(c.Request.Context(), payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "LDAP settings saved"})
}

func (ctl *SettingsController) TestLDAP(c *gin.Context) {
	var payload services.LDAPSettings
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctl.service.TestLDAP(c.Request.Context(), payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
