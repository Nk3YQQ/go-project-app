package routers

import (
	"backend/internal/auth"
	"backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func TasksRouters(router *gin.RouterGroup) {
	taskRouters := router.Group("/tasks")
	{
		taskRouters.Any("", auth.Authenticate, handlers.TaskViewSet)
		taskRouters.Any("/:id", auth.Authenticate, handlers.TaskViewSet)
	}
}
