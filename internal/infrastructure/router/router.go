package router

import (
	"fmt"
	"net/http"

	"github.com/cool9850311/StreamPlatformLite-Core/internal/application/usecase"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/domain/interface/logger"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/infrastructure/config"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/infrastructure/controller"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/infrastructure/initializer"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/infrastructure/middleware"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/infrastructure/outer_api/discord"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/infrastructure/repository"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/infrastructure/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, log logger.Logger, redisClient *redis.Client) *gin.Engine {
	r := setupRouter()
	setupMiddlewares(r)
	setupRoutes(r, db, log, redisClient)
	return r
}

func setupRouter() *gin.Engine {
	if !config.AppConfig.Server.EnableGinLog {
		r := gin.New()
		r.Use(gin.Recovery())
		return r
	}
	return gin.Default()
}

func setupMiddlewares(r *gin.Engine) {
	// Dynamic CORS configuration based on environment
	var allowedOrigins []string
	if config.AppConfig.Server.HTTPS {
		allowedOrigins = append(allowedOrigins, fmt.Sprintf("https://%s", config.AppConfig.Frontend.Domain))
	} else {
		allowedOrigins = append(allowedOrigins, fmt.Sprintf("http://%s:%d", config.AppConfig.Frontend.Domain, config.AppConfig.Frontend.Port))
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-XSRF-TOKEN"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition", "Cache-Control", "Retry-After", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"},
		AllowCredentials: true,
	}))

	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.TraceIDMiddleware())
}

func setupRoutes(r *gin.Engine, db *gorm.DB, log logger.Logger, redisClient *redis.Client) {
	// Initialize rate limiters
	initializer.InitRateLimiters()

	// Initialize repositories, use cases, and controllers
	systemSettingRepo := repository.NewPostgresSystemSettingRepository(db)
	systemSettingUseCase := usecase.NewSystemSettingUseCase(systemSettingRepo, log)
	systemSettingController := controller.NewSystemSettingController(log, systemSettingUseCase)
	discordOAuthOuterApi := discord.NewDiscordOAuthImpl(log)
	jwtGenerator := util.NewJWTLibrary()
	bcrypt := util.NewBcryptLibrary()
	stateStore := util.NewRedisStateStore(redisClient)
	discordLoginUseCase := usecase.NewDiscordLoginUseCase(systemSettingRepo, log, config.AppConfig, discordOAuthOuterApi, jwtGenerator, stateStore)
	discordOauthController := controller.NewDiscordOauthController(log, discordLoginUseCase)
	accountRepo := repository.NewPostgresAccountRepository(db)
	originAccountUseCase := usecase.NewOriginAccountUseCase(accountRepo, log, bcrypt, config.AppConfig, jwtGenerator)
	originAccountController := controller.NewOriginAccountController(log, originAccountUseCase)

	// Health check — public, no auth, used by Docker HEALTHCHECK
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	login := r.Group("/")
	{
		login.GET("/oauth/discord/init", middleware.RateLimitByIP(initializer.OAuthInitLimiter), discordOauthController.InitiateLogin)
		login.GET("/oauth/discord", discordOauthController.Callback)
		login.POST("/logout", middleware.RateLimitByIP(initializer.LogoutLimiter), discordOauthController.Logout)
	}
	r.GET("/me", middleware.OptionalJWTAuthMiddleware(log), originAccountController.GetMe)

	originAccount := r.Group("/origin-account")
	{
		originAccount.POST("/login", middleware.RateLimitByIP(initializer.LoginLimiter), originAccountController.Login)
		originAccount.POST("/create", middleware.JWTAuthMiddleware(log), originAccountController.CreateAccount)
		originAccount.PATCH("/change-password", middleware.JWTAuthMiddleware(log), middleware.RateLimitByUserID(initializer.ChangePasswordLimiter), originAccountController.ChangePassword)
		originAccount.GET("/list", middleware.JWTAuthMiddleware(log), originAccountController.GetAccountList)
		originAccount.DELETE("/delete", middleware.JWTAuthMiddleware(log), originAccountController.DeleteAccount)
	}

	systemSettings := r.Group("/system-settings")
	{
		systemSettings.GET("", middleware.JWTAuthMiddleware(log), systemSettingController.GetSetting)
		systemSettings.PATCH("", middleware.JWTAuthMiddleware(log), systemSettingController.SetSetting)
	}
}
