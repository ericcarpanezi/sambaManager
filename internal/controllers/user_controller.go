package controllers

import (
	"net/http"
	"strconv"

	"github.com/example/ag-directory-manager/internal/middleware"
	"github.com/example/ag-directory-manager/internal/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}

func (ctl *UserController) List(c *gin.Context) {
	search := c.Query("search")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	rawOU, _ := c.Get(middleware.ContextOUScopes)
	ouScopes, _ := rawOU.([]string)

	users, err := ctl.service.List(c.Request.Context(), search, ouScopes, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": users})
}
