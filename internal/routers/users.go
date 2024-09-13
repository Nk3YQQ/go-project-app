package routers

import (
	"backend/internal/auth"
	"backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func UsersRouters(router *gin.RouterGroup) {
	userRouters := router.Group("/users")
	{
		userRouters.POST("/register", handlers.RegisterUser)
		userRouters.POST("/login", handlers.LoginUser)
		userRouters.GET("/profile", auth.Authenticate, handlers.Profile)
	}
}
