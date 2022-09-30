package utils

import (
	"github.com/saiyedulbas/second/account"
	"github.com/saiyedulbas/second/dbconn"
	"gorm.io/gorm"
)

var (
	db *gorm.DB = dbconn.SetupConnection()
)

func MigrateTables() {
	db.Debug().AutoMigrate(&account.User{}, &account.Role{}) //database migration
}
