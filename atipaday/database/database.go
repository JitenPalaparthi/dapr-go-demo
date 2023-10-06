package database

import (
	"context"
	"time"

	"github.com/JitenPalaparthi/atipaday/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	MaxRetries   uint = 10
	MaxRetryTime uint = 2 // assume in seconds
)

// check retry logic with timeout
// time.After returns channel
func GetConnection(dsn string) (any, error) {
	var count uint = 1
retry:
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil && count <= MaxRetries {
		time.Sleep(time.Second * time.Duration(MaxRetryTime))
		count++
		goto retry
	}
	if err := db.AutoMigrate(&models.Tip{}); err != nil {
		return nil, err
	}
	return db, err
}

func IsExists(ctx context.Context, db any, model any, id string) (bool, error) {
	var count int64
	tx := db.(*gorm.DB).Find(model, id).Count(&count)
	//tx := vs.DB.(*gorm.DB).Model(&VehicleSpecs{}).Where("id = ?", id).Count(&count)
	if tx.Error != nil {
		return false, tx.Error
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
