package routers

import (
	"backend/internal/auth"
	"backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func ProjectRouters(router *gin.RouterGroup) {
	projectRouters := router.Group("/projects")
	{
		projectRouters.Any("", auth.Authenticate, handlers.ProjectViewSet)
		projectRouters.Any("/:id", auth.Authenticate, handlers.ProjectViewSet)
	}
}
