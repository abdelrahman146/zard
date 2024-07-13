package provider

import (
	"fmt"
	"github.com/abdelrahman146/zard/shared/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/url"
)

type GormProvider interface {
	GetDB() *gorm.DB
	Migrate(entities ...interface{})
	Close()
}

type gormProvider struct {
	db *gorm.DB
}

func InitGormProvider(address string) GormProvider {
	p := &gormProvider{}
	dsn, err := p.getDsn(address)
	if err != nil {
		logger.GetLogger().Panic("Failed to parse the database URL", logger.Field("address", address), logger.Field("error", err))
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.GetLogger().Panic("Failed to connect to the database", logger.Field("address", address), logger.Field("error", err))
	}
	logger.GetLogger().Info("Connected to the database", logger.Field("address", address))
	p.db = db
	return p
}

func (g *gormProvider) getDsn(address string) (string, error) {
	// Parse the URL
	parsedURL, err := url.Parse(address)
	if err != nil {
		return "", err
	}
	user := parsedURL.User.Username()
	password, _ := parsedURL.User.Password()
	host := parsedURL.Hostname()
	port := parsedURL.Port()
	dbName := parsedURL.Path[1:] // Skip the leading '/'

	// Extract query parameters (if any)
	queryParams := parsedURL.Query()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		host, user, password, dbName, port)
	for key, values := range queryParams {
		for _, value := range values {
			dsn += fmt.Sprintf(" %s=%s", key, value)
		}
	}
	return dsn, nil
}

func (g *gormProvider) GetDB() *gorm.DB {
	return g.db
}

func (g *gormProvider) Migrate(entities ...interface{}) {
	err := g.db.AutoMigrate(entities...)
	if err != nil {
		logger.GetLogger().Panic("Failed to migrate the database", logger.Field("error", err))
	}
	logger.GetLogger().Info("Database migration completed")
}

func (g *gormProvider) Close() {
	db, err := g.db.DB()
	if err != nil {
		logger.GetLogger().Warn("Failed to disconnect from the database", logger.Field("error", err))
		return
	}
	err = db.Close()
	if err != nil {
		logger.GetLogger().Warn("Failed to close the database connection", logger.Field("error", err))
		return
	}
	logger.GetLogger().Info("Database connection closed")
}
