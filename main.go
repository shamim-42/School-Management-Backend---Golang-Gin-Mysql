// You can edit this code!
// Click here and start typing.
package main

import (
	"github.com/gin-contrib/cors"
	"github.com/saiyedulbas/second/dbconn"
	"github.com/saiyedulbas/second/routes"
	"github.com/saiyedulbas/second/utils"
	"gorm.io/gorm"
)

var (
	db *gorm.DB = dbconn.SetupConnection()
)

func main() {
	// db connection
	defer dbconn.CloseDatabaseConnection(db)

	// db migrations
	utils.MigrateTables()

	// cors config
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Token", "Content-Type"}

	// gin router (all routes are written in routes folder)
	router := routes.Router
	router.Use(cors.New(config))
	router.Run(":8091")
}
