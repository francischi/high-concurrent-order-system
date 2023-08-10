package models

type Order struct {
	Id          int    `gorm:"type:int(10)"`
	OrderUuid     string `gorm:"type:varchar(36);not null;index:order_uuid"`
	MemberUuid    string `gorm:"type:varchar(36);not null;index:member_uuid"`
	// Price       int    `gorm:"type:int(10) unsigned;not null"`
	// Quantity    int    `gorm:"type:int(5) unsigned;not null"`
	// Status      int    `gorm:"type:tinyint;not null"`
	// Description string `gorm:"type:varchar(100);default:null"`
	CreateTime  int    `gorm:"type:int(10);not null"`
	UpdateTime  int    `gorm:"type:int(10);default:null"`
	DeleteTime  int    `gorm:"type:int(10);default:null"`
}

// add columns if you need

func (m *Order) TableName() string {
	return "orders"
}
