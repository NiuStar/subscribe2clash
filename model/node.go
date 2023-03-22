package model

import (
	"gorm.io/gorm"
)

type Node struct {
	ID      uint   `gorm:"primary_key;auto_increment"`
	Address string `gorm:"type:varchar(255);unique_index"`
	gorm.Model
}

func init() {
	db.AutoMigrate(&Node{})
}

func GetAllNodes() ([]Node, error) {
	var nodes []Node
	err := db.Find(&nodes).Error
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
