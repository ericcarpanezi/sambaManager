package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppName        string
	AppEnv         string
	Host           string
	Port           string
	DemoMode       bool
	JWTSecret      string
	AccessTTL      time.Duration
	RefreshTTL     time.Duration
	RateLimitRPS   float64
	RateLimitBurst int
	SQLitePath     string
	EncryptionKey  string

	LDAPURL                string
	LDAPBaseDN             string
	LDAPBindDN             string
	LDAPBindPassword       string
	LDAPStartTLS           bool
	LDAPInsecureSkipVerify bool
}

func Load() Config {
	return Config{
		AppName:        getEnv("APP_NAME", "AG Directory Manager"),
		AppEnv:         getEnv("APP_ENV", "development"),
		Host:           getEnv("APP_HOST", "0.0.0.0"),
		Port:           getEnv("APP_PORT", "8080"),
		DemoMode:       getEnvBool("APP_DEMO_MODE", true),
		JWTSecret:      getEnv("APP_JWT_SECRET", "change-me"),
		AccessTTL:      getEnvDuration("APP_JWT_ACCESS_TTL", "15m"),
		RefreshTTL:     getEnvDuration("APP_JWT_REFRESH_TTL", "720h"),
		RateLimitRPS:   getEnvFloat("APP_RATE_LIMIT_RPS", 5),
		RateLimitBurst: getEnvInt("APP_RATE_LIMIT_BURST", 20),
		SQLitePath:     getEnv("APP_SQLITE_PATH", "./data/agdm.db"),
		EncryptionKey:  getEnv("APP_ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef"),

		LDAPURL:                getEnv("LDAP_URL", "ldap://127.0.0.1:389"),
		LDAPBaseDN:             getEnv("LDAP_BASE_DN", "DC=dominio,DC=local"),
		LDAPBindDN:             getEnv("LDAP_BIND_DN", ""),
		LDAPBindPassword:       getEnv("LDAP_BIND_PASSWORD", ""),
		LDAPStartTLS:           getEnvBool("LDAP_START_TLS", false),
		LDAPInsecureSkipVerify: getEnvBool("LDAP_INSECURE_SKIP_VERIFY", true),
	}
}

func (c Config) ServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}

func getEnvFloat(key string, fallback float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fallback
	}
	return f
}

func getEnvDuration(key, fallback string) time.Duration {
	v := getEnv(key, fallback)
	d, err := time.ParseDuration(v)
	if err != nil {
		d, _ = time.ParseDuration(fallback)
	}
	return d
}
