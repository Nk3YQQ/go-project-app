package handlers

import (
	"backend/internal/database"
	"backend/internal/models"
	"backend/internal/utils"
	"backend/internal/validators"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type TaskHandler struct {
	DB *gorm.DB
}

func (T *TaskHandler) findTaskByID(c *gin.Context) (*models.Task, error) {
	var task models.Task
	id := c.Param("id")
	if err := T.DB.Preload("Executors").First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (T *TaskHandler) findUsersByID(ids []int) ([]models.User, error) {
	var users []models.User
	if err := T.DB.Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (T *TaskHandler) ConvertAllTasksToSchema(tasks []models.Task) []models.TaskSchema {
	var serializedTasks []models.TaskSchema

	for _, task := range tasks {
		serializedTasks = append(serializedTasks, task.ToSchema())
	}

	return serializedTasks
}

func (T *TaskHandler) CreateTask(c *gin.Context) {
	var input models.TaskCreateSchema
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := T.findUsersByID(input.Executors)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't find users"})
		return
	}

	var project models.Project
	if err := T.DB.First(&project, "id = ?", input.ProjectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't find project"})
		return
	}

	var ParsedDeadline time.Time
	if err := utils.ParseDateToTime(input.Deadline, &ParsedDeadline); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect type of deadline"})
		return
	}

	if err := validators.ValidateDates(nil, &ParsedDeadline); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Deadline:    ParsedDeadline,
		Status:      input.Status,
		ProjectID:   project.ID,
		Executors:   users,
	}

	T.DB.Create(&task)

	c.JSON(http.StatusCreated, task.ToSchema())
}

func (T *TaskHandler) ReadTasks(c *gin.Context) {
	var tasks []models.Task

	if err := T.DB.Preload("Executors").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't find projects"})
		return
	}

	c.JSON(http.StatusOK, T.ConvertAllTasksToSchema(tasks))
}

func (T *TaskHandler) ReadTask(c *gin.Context) {
	task, err := T.findTaskByID(c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving task"})
		}
		return
	}
	c.JSON(http.StatusOK, task.ToSchema())
}

func (T *TaskHandler) UpdateTask(c *gin.Context) {
	task, err := T.findTaskByID(c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving task"})
		}
		return
	}

	var input models.TaskUpdateSchema
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ParsedDeadline time.Time
	if err := utils.ParseDateToTime(input.Deadline, &ParsedDeadline); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect type of deadline"})
		return
	}

	if err := validators.ValidateDates(nil, &ParsedDeadline); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := T.findUsersByID(input.Executors)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't find users"})
		return
	}

	task.Title = input.Title
	task.Description = input.Description
	task.Deadline = ParsedDeadline
	if input.Status != "" {
		task.Status = input.Status
	}
	if input.Status != "" {
		task.ProjectID = uint(input.ProjectID)
	}
	task.Executors = users

	T.DB.Save(&task)

	c.JSON(http.StatusOK, task.ToSchema())
}

func (T *TaskHandler) DeleteTask(c *gin.Context) {
	task, err := T.findTaskByID(c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving task"})
		}
		return
	}

	database.DB.Delete(&task)
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func TaskViewSet(c *gin.Context) {
	taskHandler := TaskHandler{DB: database.DB}

	switch c.Request.Method {
	case "GET":
		id := c.Param("id")
		if id != "" {
			taskHandler.ReadTask(c)
		} else {
			taskHandler.ReadTasks(c)
		}
	case "POST":
		taskHandler.CreateTask(c)
	case "PUT":
		taskHandler.UpdateTask(c)
	case "DELETE":
		taskHandler.DeleteTask(c)
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method now allowed"})
	}
}
