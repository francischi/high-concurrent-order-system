package implement

import (
	// "fmt"
	"gorm.io/gorm"
	"golang/pkg/base"
	"golang/pkg/helpers"
	"golang/pkg/repos/models"
)

type OrderRepo struct{
	base.Repository
	DBconn *gorm.DB
	ConnPool *helpers.ConnPool
}

func NewOrderRepo(DBconn *gorm.DB,ConnPool *helpers.ConnPool) *OrderRepo{
	return &OrderRepo{
		ConnPool:ConnPool,
		DBconn:DBconn,
	}
}

func(o *OrderRepo) AddOrder(memberId string  , orderProducts map[string]int)(err error){
	var order models.Order
	order.OrderUuid = helpers.CreateUuid()
	order.MemberUuid = memberId
	order.CreateTime = helpers.GetTimeStamp()
	if err := o.DBconn.Create(&order).Error;err!=nil{
		return o.SystemError(err.Error())
	}
	
	return nil
}