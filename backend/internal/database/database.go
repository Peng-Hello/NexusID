package database

import (
	"fmt"
	"time"

	"github.com/nexus-id/backend/internal/config"
	"github.com/nexus-id/backend/internal/models"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB holds the database connection
var DB *gorm.DB

// Init initializes the database connection
func Init(cfg *config.Config, log *zap.Logger) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DatabaseName,
	)

	// Configure GORM logger
	gormLogger := logger.Default
	if cfg.Server.Mode == "production" || cfg.Server.Mode == "release" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	log.Info("Database connected successfully",
		zap.String("host", cfg.Database.Host),
		zap.Int("port", cfg.Database.Port),
	)

	// Auto-migrate models (for development)
	if cfg.Server.Mode == "debug" {
		if err := autoMigrate(); err != nil {
			return fmt.Errorf("failed to auto-migrate: %w", err)
		}
	}

	// Seed default tenant if not exists
	var count int64
	DB.Model(&models.Tenant{}).Count(&count)
	if count == 0 {
		tenant := models.Tenant{
			Name:        "Default Tenant",
			Slug:        "default",
			Description: "The default system tenant",
			IsActive:    true,
		}
		if err := DB.Create(&tenant).Error; err != nil {
			log.Warn("Failed to seed default tenant", zap.Error(err))
		} else {
			log.Info("Seeded default tenant", zap.Int64("id", tenant.ID))
		}
	}

	return nil
}

// autoMigrate runs auto-migration for all models
func autoMigrate() error {
	return DB.AutoMigrate(
		&models.Tenant{},
		&models.OIDCClient{},
		&models.User{},
		&models.Role{},
		&models.UserRole{},
		&models.RefreshToken{},
		&models.AuthorizationCode{},
		&models.Session{},
	)
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return DB
}
