package controllers

import (
	"errors"
	"net/http"

	"github.com/example/ag-directory-manager/internal/dto"
	"github.com/example/ag-directory-manager/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (ctl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIError{Error: err.Error()})
		return
	}

	res, err := ctl.service.Login(c.Request.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrUnauthorized) {
			status = http.StatusUnauthorized
		} else if errors.Is(err, services.ErrForbidden) {
			status = http.StatusForbidden
		}
		c.JSON(status, dto.APIError{Error: err.Error()})
		return
	}

	c.SetCookie("csrf_token", "agdm-static-csrf", 3600, "/", "", false, true)
	c.JSON(http.StatusOK, res)
}

func (ctl *AuthController) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIError{Error: err.Error()})
		return
	}

	res, err := ctl.service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.APIError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (ctl *AuthController) Logout(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIError{Error: err.Error()})
		return
	}
	uidAny, _ := c.Get("userID")
	uid, _ := uidAny.(int64)

	if err := ctl.service.Logout(c.Request.Context(), uid, req.RefreshToken); err != nil {
		c.JSON(http.StatusUnauthorized, dto.APIError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.APIMessage{Message: "logout successful"})
}
