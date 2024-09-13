package handlers

import (
	"backend/internal/database"
	"backend/internal/models"
	"backend/internal/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type UserHandler struct {
	DB *gorm.DB
}

func (u *UserHandler) findProjectByID(c *gin.Context) (*models.Project, error) {
	var project models.Project
	id := c.Param("id")
	if err := u.DB.Preload("Executors").First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (u *UserHandler) findUsersByID(ids []int) ([]models.User, error) {
	var users []models.User
	if err := u.DB.Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserHandler) ConvertAllProjectsToSchema(projects []models.Project) []models.ProjectSchema {
	var serializedProjects []models.ProjectSchema

	for _, project := range projects {
		serializedProjects = append(serializedProjects, project.ToSchema())
	}

	return serializedProjects
}

func (u *UserHandler) CreateProject(c *gin.Context) {
	var input models.ProjectCreateSchema
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := u.findUsersByID(input.Executors)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't find users"})
		return
	}

	var ParsedStartedAt time.Time
	if err := utils.ParseDateToTime(input.StartedAt, &ParsedStartedAt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect type of started_at"})
		return
	}

	var ParsedDeadline time.Time
	if err := utils.ParseDateToTime(input.Deadline, &ParsedDeadline); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect type of deadline"})
		return
	}

	currentTime := time.Now()

	if currentTime.After(ParsedStartedAt) || currentTime.After(ParsedDeadline) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "current date cannot be later than started_at or deadline"})
		return
	}

	if ParsedDeadline.Before(ParsedStartedAt) || ParsedDeadline.Equal(ParsedStartedAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deadline must be after started_at"})
		return
	}

	project := models.Project{
		Title:       input.Title,
		Description: input.Description,
		StartedAt:   ParsedStartedAt,
		Deadline:    ParsedDeadline,
		Status:      input.Status,
		Executors:   users,
		Tasks:       input.Tasks,
	}

	database.DB.Create(&project)

	c.JSON(http.StatusCreated, project.ToSchema())
}

func (u *UserHandler) ReadProjects(c *gin.Context) {
	var projects []models.Project
	if err := u.DB.Preload("Executors").Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't find projects"})
		return
	}

	c.JSON(http.StatusOK, u.ConvertAllProjectsToSchema(projects))
}

func (u *UserHandler) ReadProject(c *gin.Context) {
	project, err := u.findProjectByID(c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving project"})
		}
		return
	}
	c.JSON(http.StatusOK, project.ToSchema())
}

func (u *UserHandler) UpdateProject(c *gin.Context) {
	project, err := u.findProjectByID(c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving project"})
		}
		return
	}

	var input models.ProjectUpdateSchema
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ParsedStartedAt time.Time
	if err := utils.ParseDateToTime(input.StartedAt, &ParsedStartedAt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect type of started_at"})
		return
	}

	var ParsedDeadline time.Time
	if err := utils.ParseDateToTime(input.Deadline, &ParsedDeadline); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect type of deadline"})
		return
	}

	currentTime := time.Now()

	if currentTime.After(ParsedStartedAt) || currentTime.After(ParsedDeadline) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "current date cannot be later than started_at or deadline"})
		return
	}

	if ParsedDeadline.Before(ParsedStartedAt) || ParsedDeadline.Equal(ParsedStartedAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deadline must be after started_at"})
		return
	}

	users, err := u.findUsersByID(input.Executors)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't find users"})
		return
	}

	project.Title = input.Title
	project.Description = input.Description
	project.StartedAt = ParsedStartedAt
	project.Deadline = ParsedDeadline
	if input.Status != "" {
		project.Status = input.Status
	}
	project.Executors = users
	project.Tasks = input.Tasks

	u.DB.Save(&project)

	c.JSON(http.StatusOK, project.ToSchema())
}

func (u *UserHandler) DeleteProject(c *gin.Context) {
	project, err := u.findProjectByID(c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving project"})
		}
		return
	}
	c.JSON(http.StatusOK, project.ToSchema())

	u.DB.Delete(&project)
	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

func ProjectViewSet(c *gin.Context) {
	userHandler := UserHandler{DB: database.DB}
	switch c.Request.Method {
	case "GET":
		id := c.Param("id")
		if id != "" {
			userHandler.ReadProject(c)
		} else {
			userHandler.ReadProjects(c)
		}
	case "POST":
		userHandler.CreateProject(c)
	case "PUT":
		userHandler.UpdateProject(c)
	case "DELETE":
		userHandler.DeleteProject(c)
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method now allowed"})
	}
}
