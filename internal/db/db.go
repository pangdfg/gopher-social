package db

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGorm(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string, debug bool) (*gorm.DB, error) {
	// Parse max idle time
	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}

	// GORM config
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // silent by default
	}

	if debug {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(addr), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxIdleTime(duration)
	sqlDB.SetConnMaxLifetime(0)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
