package migrations

import (
	"golang/pkg/repos/models"
	"gorm.io/gorm"
)

func Product(db *gorm.DB){
	var ProductModel models.ProductModel
	db.AutoMigrate(&ProductModel)
}