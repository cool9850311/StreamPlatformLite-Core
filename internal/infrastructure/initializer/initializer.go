package initializer

import (
	"log"

	domainLogger "github.com/cool9850311/StreamPlatformLite-Core/internal/domain/interface/logger"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/infrastructure/config"
	infraLogger "github.com/cool9850311/StreamPlatformLite-Core/internal/infrastructure/logger"
	migrations "github.com/cool9850311/StreamPlatformLite-Core/migrations"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GormDB *gorm.DB
var Log domainLogger.Logger
var RedisClient *redis.Client

func InitLog() {
	var err error
	logLevel := config.AppConfig.Server.LogLevel
	if logLevel == "" {
		logLevel = "INFO"
	}
	Log, err = infraLogger.NewLogger("application.log", logLevel)
	if err != nil {
		panic(err)
	}
}

func InitConfig() {
	config.LoadConfig()
}

func InitSchema() {
	if !config.AppConfig.PostgreSQL.AutoMigrateSchema {
		log.Println("[schema] Auto-migrate disabled, skipping")
		return
	}
	d, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatalf("failed to create migration source: %v", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, config.AppConfig.PostgreSQL.DSN)
	if err != nil {
		log.Fatalf("failed to create migrator: %v", err)
	}
	defer m.Close()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("schema migration failed: %v", err)
	}
	log.Println("[schema] Schema migration done.")
}

func InitPostgresClient() {
	dsn := config.AppConfig.PostgreSQL.DSN
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}
	GormDB = db
	log.Println("Connected to PostgreSQL")
}

func InitRedisClient() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.Redis.URI,
		Password: "",
		DB:       0,
	})
}
