package services

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/example/ag-directory-manager/internal/config"
	"github.com/example/ag-directory-manager/internal/ldap"
)

type SettingsService struct {
	db  *sql.DB
	cfg config.Config
}

func NewSettingsService(db *sql.DB, cfg config.Config) *SettingsService {
	return &SettingsService{db: db, cfg: cfg}
}

type LDAPSettings struct {
	ServerURL  string `json:"serverUrl"`
	BaseDN     string `json:"baseDn"`
	BindDN     string `json:"bindDn"`
	BindSecret string `json:"bindSecret"`
	StartTLS   bool   `json:"startTls"`
	SkipVerify bool   `json:"skipVerify"`
}

func (s *SettingsService) GetLDAP(ctx context.Context) (LDAPSettings, error) {
	keys := []string{"ldap.url", "ldap.base_dn", "ldap.bind_dn", "ldap.bind_secret", "ldap.start_tls", "ldap.skip_verify"}
	values := map[string]string{}

	for _, k := range keys {
		var v string
		err := s.db.QueryRowContext(ctx, `SELECT value FROM settings WHERE key = ?`, k).Scan(&v)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			return LDAPSettings{}, err
		}
		values[k] = v
	}

	bindSecret := decodeValue(values["ldap.bind_secret"])
	if bindSecret == "" {
		bindSecret = s.cfg.LDAPBindPassword
	}

	return LDAPSettings{
		ServerURL:  fallback(values["ldap.url"], s.cfg.LDAPURL),
		BaseDN:     fallback(values["ldap.base_dn"], s.cfg.LDAPBaseDN),
		BindDN:     fallback(values["ldap.bind_dn"], s.cfg.LDAPBindDN),
		BindSecret: bindSecret,
		StartTLS:   strings.EqualFold(values["ldap.start_tls"], "true") || (values["ldap.start_tls"] == "" && s.cfg.LDAPStartTLS),
		SkipVerify: !strings.EqualFold(values["ldap.skip_verify"], "false") && (values["ldap.skip_verify"] != "" || s.cfg.LDAPInsecureSkipVerify),
	}, nil
}

func (s *SettingsService) SaveLDAP(ctx context.Context, cfg LDAPSettings) error {
	updates := map[string]string{
		"ldap.url":         cfg.ServerURL,
		"ldap.base_dn":     cfg.BaseDN,
		"ldap.bind_dn":     cfg.BindDN,
		"ldap.bind_secret": encodeValue(cfg.BindSecret),
		"ldap.start_tls":   fmt.Sprintf("%v", cfg.StartTLS),
		"ldap.skip_verify": fmt.Sprintf("%v", cfg.SkipVerify),
	}
	for k, v := range updates {
		if _, err := s.db.ExecContext(ctx, `
			INSERT INTO settings(key, value, updated_at)
			VALUES (?, ?, CURRENT_TIMESTAMP)
			ON CONFLICT(key) DO UPDATE SET value=excluded.value, updated_at=CURRENT_TIMESTAMP
		`, k, v); err != nil {
			return err
		}
	}
	return nil
}

func (s *SettingsService) TestLDAP(ctx context.Context, cfg LDAPSettings) error {
	client := ldap.NewRealClient(ldap.Config{
		URL:                cfg.ServerURL,
		BaseDN:             cfg.BaseDN,
		BindDN:             cfg.BindDN,
		BindPassword:       cfg.BindSecret,
		StartTLS:           cfg.StartTLS,
		InsecureSkipVerify: cfg.SkipVerify,
	})
	return client.TestConnection(ctx)
}

func encodeValue(value string) string {
	if value == "" {
		return ""
	}
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func decodeValue(value string) string {
	if value == "" {
		return ""
	}
	raw, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return ""
	}
	return string(raw)
}

func fallback(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
