package api

import (
	"database/sql"
	"errors"
	"time"

	"github.com/example/ag-directory-manager/internal/audit"
	"github.com/example/ag-directory-manager/internal/config"
	"github.com/example/ag-directory-manager/internal/controllers"
	"github.com/example/ag-directory-manager/internal/ldap"
	"github.com/example/ag-directory-manager/internal/middleware"
	"github.com/example/ag-directory-manager/internal/models"
	"github.com/example/ag-directory-manager/internal/repositories"
	"github.com/example/ag-directory-manager/internal/services"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config, db *sql.DB) (*gin.Engine, error) {
	if cfg.JWTSecret == "" || cfg.JWTSecret == "change-me" {
		return nil, errors.New("APP_JWT_SECRET must be configured")
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst).Middleware())

	userRepo := repositories.NewUserRepository(db)
	auditRepo := repositories.NewAuditRepository(db)
	tokenRepo := repositories.NewTokenRepository(db)
	dashRepo := repositories.NewDashboardRepository(db)

	var directory ldap.Client
	if cfg.DemoMode {
		directory = ldap.NewDemoClient()
	} else {
		directory = ldap.NewRealClient(ldap.Config{
			URL:                cfg.LDAPURL,
			BaseDN:             cfg.LDAPBaseDN,
			BindDN:             cfg.LDAPBindDN,
			BindPassword:       cfg.LDAPBindPassword,
			StartTLS:           cfg.LDAPStartTLS,
			InsecureSkipVerify: cfg.LDAPInsecureSkipVerify,
		})
	}

	tokenManager := services.NewTokenManager(cfg.JWTSecret, cfg.AccessTTL, cfg.RefreshTTL)
	authService := services.NewAuthService(userRepo, tokenRepo, directory, tokenManager)
	dashboardService := services.NewDashboardService(directory, dashRepo)
	userService := services.NewUserService(directory)
	settingsService := services.NewSettingsService(db, cfg)
	auditService := audit.NewService(auditRepo)

	healthCtl := controllers.NewHealthController()
	authCtl := controllers.NewAuthController(authService)
	dashCtl := controllers.NewDashboardController(dashboardService)
	usersCtl := controllers.NewUserController(userService)
	auditCtl := controllers.NewAuditController(auditService)
	settingsCtl := controllers.NewSettingsController(settingsService)
	events := newEventHub()

	r.GET("/healthz", healthCtl.Check)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/login", authCtl.Login)
		v1.POST("/auth/refresh", authCtl.Refresh)

		secured := v1.Group("")
		secured.Use(middleware.AuthRequired(authService))
		secured.Use(middleware.CSRFMiddleware())
		secured.Use(func(c *gin.Context) {
			start := time.Now()
			c.Next()

			usernameAny, _ := c.Get(middleware.ContextUsername)
			username, _ := usernameAny.(string)
			auditService.Record(c.Request.Context(), models.AuditLog{
				ActorUsername: username,
				IPAddress:     c.ClientIP(),
				UserAgent:     c.GetHeader("User-Agent"),
				Operation:     c.Request.Method,
				ObjectType:    "api",
				ObjectID:      c.FullPath(),
				Result:        httpStatusResult(c.Writer.Status()),
				DurationMS:    time.Since(start).Milliseconds(),
			})
		})
		{
			secured.POST("/auth/logout", authCtl.Logout)
			secured.GET("/dashboard", middleware.RequirePermission("user.view"), dashCtl.Snapshot)
			secured.GET("/users", middleware.RequirePermission("user.view"), usersCtl.List)
			secured.GET("/audit/logs", middleware.RequirePermission("audit.view"), auditCtl.List)

			secured.GET("/settings/ldap", middleware.RequirePermission("settings.manage"), settingsCtl.GetLDAP)
			secured.PUT("/settings/ldap", middleware.RequirePermission("settings.manage"), settingsCtl.SaveLDAP)
			secured.POST("/settings/ldap/test", middleware.RequirePermission("settings.manage"), settingsCtl.TestLDAP)

			secured.GET("/ws/events", events.Handle)
		}
	}

	r.GET("/docs/swagger", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Swagger integration placeholder. Generate OpenAPI via docs/openapi.yaml"})
	})

	return r, nil
}

func httpStatusResult(status int) string {
	if status >= 200 && status < 400 {
		return "success"
	}
	return "error"
}
