package migrations

import (
	"golang/pkg/repos/models"
	"gorm.io/gorm"
)

func Order(db *gorm.DB){
	var OrderModel models.OrderModel
	db.AutoMigrate(&OrderModel)
}