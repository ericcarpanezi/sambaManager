package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/example/ag-directory-manager/internal/dto"
	"github.com/example/ag-directory-manager/internal/ldap"
	"github.com/example/ag-directory-manager/internal/repositories"
)

var (
	ErrUnauthorized = errors.New("invalid credentials")
	ErrForbidden    = errors.New("forbidden")
)

type AuthService struct {
	users  *repositories.UserRepository
	tokens *repositories.TokenRepository
	dir    ldap.Client
	jwt    *TokenManager
}

func NewAuthService(users *repositories.UserRepository, tokens *repositories.TokenRepository, dir ldap.Client, jwt *TokenManager) *AuthService {
	return &AuthService{users: users, tokens: tokens, dir: dir, jwt: jwt}
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	if err := s.dir.Authenticate(ctx, req.Username, req.Password); err != nil {
		return dto.LoginResponse{}, ErrUnauthorized
	}

	user, err := s.users.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return dto.LoginResponse{}, ErrForbidden
		}
		return dto.LoginResponse{}, err
	}

	if !user.IsActive {
		return dto.LoginResponse{}, ErrForbidden
	}

	if user.BlockedUntil != nil && user.BlockedUntil.After(time.Now().UTC()) {
		return dto.LoginResponse{}, fmt.Errorf("account temporarily blocked")
	}

	permissions, err := s.users.ListPermissions(ctx, user.ID)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	ouScopes, err := s.users.ListOUScopes(ctx, user.ID)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	access, accessExp, err := s.jwt.NewAccessToken(user.ID, user.Username, permissions, ouScopes)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	rawRefresh, refreshHash, refreshExp, err := s.jwt.NewRefreshToken(user.ID)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	if err := s.tokens.StoreRefreshToken(ctx, user.ID, refreshHash, refreshExp); err != nil {
		return dto.LoginResponse{}, err
	}

	_ = s.users.UpdateLastLogin(ctx, user.ID)

	return dto.LoginResponse{
		User:        user,
		Permissions: permissions,
		OUScopes:    ouScopes,
		Tokens: dto.TokenPair{
			AccessToken:  access,
			RefreshToken: rawRefresh,
			TokenType:    "Bearer",
			ExpiresIn:    int64(accessExp.Sub(time.Now().UTC()).Seconds()),
		},
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, rawRefresh string) (dto.TokenPair, error) {
	claims, refreshHash, err := s.jwt.ParseRefreshToken(rawRefresh)
	if err != nil {
		return dto.TokenPair{}, ErrUnauthorized
	}

	if err := s.tokens.ValidateRefreshToken(ctx, claims.UserID, refreshHash); err != nil {
		return dto.TokenPair{}, ErrUnauthorized
	}

	user, err := s.users.GetByID(ctx, claims.UserID)
	if err != nil {
		return dto.TokenPair{}, ErrUnauthorized
	}

	permissions, err := s.users.ListPermissions(ctx, user.ID)
	if err != nil {
		return dto.TokenPair{}, err
	}
	ouScopes, err := s.users.ListOUScopes(ctx, user.ID)
	if err != nil {
		return dto.TokenPair{}, err
	}

	access, accessExp, err := s.jwt.NewAccessToken(user.ID, user.Username, permissions, ouScopes)
	if err != nil {
		return dto.TokenPair{}, err
	}

	return dto.TokenPair{
		AccessToken:  access,
		RefreshToken: rawRefresh,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessExp.Sub(time.Now().UTC()).Seconds()),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, userID int64, refreshToken string) error {
	claims, tokenHash, err := s.jwt.ParseRefreshToken(refreshToken)
	if err != nil {
		return ErrUnauthorized
	}
	if claims.UserID != userID {
		return ErrUnauthorized
	}
	if err := s.tokens.RevokeRefreshToken(ctx, userID, tokenHash); err != nil {
		return ErrUnauthorized
	}
	return nil
}

func (s *AuthService) ParseAccessToken(raw string) (*AccessClaims, error) {
	return s.jwt.ParseAccessToken(raw)
}

func ParseUserID(subject string) (int64, error) {
	return strconv.ParseInt(subject, 10, 64)
}
