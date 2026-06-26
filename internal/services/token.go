package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewTokenManager(secret string, accessTTL, refreshTTL time.Duration) *TokenManager {
	return &TokenManager{secret: []byte(secret), accessTTL: accessTTL, refreshTTL: refreshTTL}
}

type AccessClaims struct {
	UserID      int64    `json:"uid"`
	Username    string   `json:"username"`
	Permissions []string `json:"permissions"`
	OUScopes    []string `json:"ouScopes"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID int64 `json:"uid"`
	jwt.RegisteredClaims
}

func (m *TokenManager) NewAccessToken(userID int64, username string, permissions, ouScopes []string) (string, time.Time, error) {
	now := time.Now().UTC()
	expires := now.Add(m.accessTTL)
	claims := AccessClaims{
		UserID:      userID,
		Username:    username,
		Permissions: permissions,
		OUScopes:    ouScopes,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ag-directory-manager",
			Subject:   strconv.FormatInt(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expires),
			ID:        uuid.NewString(),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString(m.secret)
	return token, expires, err
}

func (m *TokenManager) NewRefreshToken(userID int64) (string, string, time.Time, error) {
	now := time.Now().UTC()
	expires := now.Add(m.refreshTTL)
	claims := RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ag-directory-manager",
			Subject:   strconv.FormatInt(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expires),
			ID:        uuid.NewString(),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	raw, err := t.SignedString(m.secret)
	if err != nil {
		return "", "", time.Time{}, err
	}
	h := sha256.Sum256([]byte(raw))
	return raw, hex.EncodeToString(h[:]), expires, nil
}

func (m *TokenManager) ParseAccessToken(raw string) (*AccessClaims, error) {
	parsed, err := jwt.ParseWithClaims(raw, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*AccessClaims)
	if !ok || !parsed.Valid {
		return nil, fmt.Errorf("invalid access token")
	}
	return claims, nil
}

func (m *TokenManager) ParseRefreshToken(raw string) (*RefreshClaims, string, error) {
	parsed, err := jwt.ParseWithClaims(raw, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, "", err
	}
	claims, ok := parsed.Claims.(*RefreshClaims)
	if !ok || !parsed.Valid {
		return nil, "", fmt.Errorf("invalid refresh token")
	}
	h := sha256.Sum256([]byte(raw))
	return claims, hex.EncodeToString(h[:]), nil
}
