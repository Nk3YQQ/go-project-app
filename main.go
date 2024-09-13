package main

import (
	"backend/internal/database"
	"backend/internal/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	router := gin.Default()

	APIRouter := router.Group("/api/v1")

	routers.ProjectRouters(APIRouter)
	routers.UsersRouters(APIRouter)
	routers.TasksRouters(APIRouter)

	router.Run()
}
