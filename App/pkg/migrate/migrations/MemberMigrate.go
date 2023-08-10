package migrations

import (
	"golang/pkg/repos/models"
	"gorm.io/gorm"
)

func Member(db *gorm.DB){
	var MemberModel models.Member
	db.AutoMigrate(&MemberModel)
}