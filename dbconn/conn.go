package dbconn

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupConnection() *gorm.DB {
	dsn := "shamim:PASSWORD!@tcp(127.0.0.1:3306)/go_gorm_mysql?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	_ = err
	return db
}

func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic("Failed to close connection from database")
	}
	dbSQL.Close()
}
