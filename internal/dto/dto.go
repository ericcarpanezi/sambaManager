package dto

import "github.com/example/ag-directory-manager/internal/models"

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int64  `json:"expiresIn"`
}

type LoginResponse struct {
	User        models.AppUser `json:"user"`
	Permissions []string       `json:"permissions"`
	OUScopes    []string       `json:"ouScopes"`
	Tokens      TokenPair      `json:"tokens"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type APIError struct {
	Error string `json:"error"`
}

type APIMessage struct {
	Message string `json:"message"`
}

type DashboardResponse struct {
	Snapshot models.DashboardSnapshot `json:"snapshot"`
}
