package migrations

import (
	"golang/pkg/repos/models"
	"golang/pkg/helpers"
	"gorm.io/gorm"
)

func Product(db *gorm.DB){

	var ProductModel models.Product
	if !db.Migrator().HasTable(ProductModel.TableName()){
		db.AutoMigrate(&ProductModel)
		seedProducts(db)
	}
}

func seedProducts(db *gorm.DB)error{
	products := []*models.Product{
		{
			ProductUuid: "a9ca26df-9c8a-4d41-85b5-5afd605e0e2f",
			Name: "product1",
			Price: 2000,
			Available: 15,
			CreateTime: helpers.GetTimeStamp(),	
		},
		{
			ProductUuid: "4e0c3b77-ba50-477b-85d1-b99bb6129115",
			Name: "product2",
			Price: 1500,
			Available: 10,
			CreateTime: helpers.GetTimeStamp(),	
		},
	}

	return db.Create(&products).Error
}

func nilInt() *int {
	return nil
}