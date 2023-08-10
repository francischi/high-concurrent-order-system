package implement

import (
	// "fmt"
	// "errors"
	"gorm.io/gorm"
	"productConsumer/pkg/base"
	"productConsumer/pkg/repos/models"
)

type ProductRepo struct {
	base.Repository
	DBconn *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepo{
	return &ProductRepo{
		DBconn:db,
	}
}

func (p *ProductRepo)ReduceProducts(productIds map[string]int)(err error){
	var productModel models.ProductModel
	var productModels []models.ProductModel
    productIDs := make([]string, 0, len(productIds))

    for productId := range productIds {
        productIDs = append(productIDs, productId)
    }


	tx := p.DBconn.Begin()

	if err := tx.Table(productModel.TableName()).Where("product_uuid IN (?)", productIDs).Set("gorm:query_option", "FOR UPDATE").Find(&productModels).Error; err != nil {
        tx.Rollback()
        return p.SystemError(err.Error())
    }

	if len (productModels) != len(productIds){
		tx.Rollback()
        return p.SystemError("lock error")
	}

	for _, productModel := range productModels {
        quantityToReduce, exists := productIds[productModel.ProductUuid]
        if !exists || productModel.Available < quantityToReduce {
            tx.Rollback()
            return p.SystemError("product shortage")
        }

        productModel.Available -= quantityToReduce
        if err := tx.Save(&productModel).Error; err != nil {
            tx.Rollback()
            return p.SystemError("product shortage")
        }
    }

	tx.Commit()
	return nil
}
