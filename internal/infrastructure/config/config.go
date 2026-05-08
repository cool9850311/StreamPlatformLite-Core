package config

import (
	"log"
	"os"
	"strconv"

	dtoConfig "github.com/cool9850311/StreamPlatformLite-Core/internal/application/dto/config"
	"github.com/joho/godotenv"
)

var AppConfig dtoConfig.Config

// getEnvAsBool reads an environment variable as a boolean with a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Printf("Invalid value for %s: %s, using default: %t", key, err, defaultValue)
		return defaultValue
	}
	return value
}

// getEnvAsInt64 reads an environment variable as int64 with a default value
func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		log.Printf("Invalid value for %s: %s, using default: %d", key, err, defaultValue)
		return defaultValue
	}
	return value
}

func LoadConfig() {
	// Load .env file from current working directory
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Read environment variables
	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatalf("Invalid SERVER_PORT: %s", err)
	}
	AppConfig.Server.Port = port
	dsn := os.Getenv("POSTGRESQL_DSN")
	if dsn == "" {
		log.Fatal("POSTGRESQL_DSN is required")
	}
	AppConfig.PostgreSQL.DSN = dsn
	AppConfig.PostgreSQL.AutoMigrateSchema = getEnvAsBool("SCHEMA_AUTO_MIGRATE", true)
	AppConfig.JWT.SecretKey = os.Getenv("APP_SECRET_KEY")
	AppConfig.Discord.ClientID = os.Getenv("DISCORD_CLIENT_ID")
	AppConfig.Discord.ClientSecret = os.Getenv("DISCORD_CLIENT_SECRET")
	AppConfig.Discord.AdminID = os.Getenv("DISCORD_ADMIN_ID")
	AppConfig.Discord.GuildID = os.Getenv("DISCORD_GUILD_ID")
	AppConfig.Server.Domain = os.Getenv("DOMAIN")
	AppConfig.Frontend.Domain = os.Getenv("FRONTEND_DOMAIN")
	AppConfig.Frontend.Port = int(getEnvAsInt64("FRONTEND_PORT", 3000))
	loginPath := os.Getenv("FRONTEND_LOGIN_PATH")
	if loginPath == "" {
		loginPath = "/stream"
	}
	AppConfig.Frontend.LoginPath = loginPath
	AppConfig.Redis.URI = os.Getenv("REDIS_URI")
	AppConfig.Server.EnableGinLog, err = strconv.ParseBool(os.Getenv("ENABLE_GIN_LOG"))
	if err != nil {
		log.Printf("Invalid ENABLE_GIN_LOG: %s", err)
	}
	AppConfig.Server.HTTPS, err = strconv.ParseBool(os.Getenv("HTTPS"))
	if err != nil {
		log.Printf("Invalid HTTPS: %s", err)
	}

	// Load log level, default to INFO if not set
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}
	AppConfig.Server.LogLevel = logLevel

	// Load Rate Limiting configuration
	AppConfig.RateLimit.Enabled = getEnvAsBool("RATE_LIMIT_ENABLED", true)
	AppConfig.RateLimit.LoginPerMinute = getEnvAsInt64("RATE_LIMIT_LOGIN_PER_MINUTE", 5)
	AppConfig.RateLimit.OAuthInitPerMinute = getEnvAsInt64("RATE_LIMIT_OAUTH_INIT_PER_MINUTE", 5)
	AppConfig.RateLimit.LogoutPerMinute = getEnvAsInt64("RATE_LIMIT_LOGOUT_PER_MINUTE", 10)
	AppConfig.RateLimit.ChangePasswordPerHour = getEnvAsInt64("RATE_LIMIT_CHANGE_PASSWORD_PER_HOUR", 10)
}
