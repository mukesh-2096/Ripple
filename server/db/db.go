package db

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// RequestLog represents the GORM schema for logging API requests
type RequestLog struct {
	gorm.Model
	Method       string  `json:"method"`
	URL          string  `json:"url"`
	StatusCode   int     `json:"status_code"`
	ResponseTime float64 `json:"response_time"` // in milliseconds
	Profile      string  `json:"profile"`       // e.g. 5G, 4G, etc.
}

func InitDB(dsn string) error {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		DB = nil
		return err
	}

	// Optimize connection pool for load handling (Mandate 4)
	sqlDB, err := DB.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(50)
		sqlDB.SetMaxOpenConns(200)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	return DB.AutoMigrate(&RequestLog{}, &User{}, &ContactMessage{})
}
