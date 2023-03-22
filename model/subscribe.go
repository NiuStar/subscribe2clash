package model

import (
	"fmt"
	"gorm.io/gorm"
)

type Subscribe struct {
	ID           uint   `gorm:"primary_key;auto_increment"`
	ShortCode    string `gorm:"unique_index;not null"`
	SubscribeURL string `gorm:"not null"`
	gorm.Model
}

func init() {
	err := db.AutoMigrate(&Subscribe{})
	fmt.Println("err", err)
}

func GetSubscribeByShortCode(shortCode string) (*Subscribe, error) {
	var subscribe Subscribe
	err := db.Where("short_code = ?", shortCode).First(&subscribe).Error
	if err != nil {
		return nil, err
	}
	return &subscribe, nil
}
