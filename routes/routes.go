package routes

import (
	"github.com/gin-gonic/gin"
	account "github.com/saiyedulbas/second/account"
)

var Router = gin.Default()

func init() {
	authRoutes := Router.Group("api/auth")
	{
		authRoutes.POST("/role", account.CreateRole)
		authRoutes.GET("/role", account.ListRole)
		authRoutes.GET("/role/:id", account.GetRole)
		authRoutes.PATCH("/role/:id", account.UpdateRole)
		authRoutes.POST("/signup", account.Signup)
		authRoutes.POST("/signin", account.Signin)
		authRoutes.PATCH("/update-user/:id", account.UpdateUser)

		// authRoutes.GET("/convertUrl", account.ConvertUrl)

	}
}
