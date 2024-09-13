package handlers

import (
	"backend/internal/auth"
	"backend/internal/database"
	"backend/internal/models"
	"backend/internal/utils"
	"backend/internal/validators"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func RegisterUser(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method now allowed"})
		return
	}

	var input models.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RaiseBadRequestError(c, err)
		return
	}

	validationErrors := validators.ValidateRegisterForm(input.FirstName, input.LastName, input.Email, input.Password, input.PasswordConfirm)

	for k := range validationErrors {
		if validationErrors[k] != "" {
			c.JSON(http.StatusBadRequest, validators.ErrorResponse{Details: validationErrors})
			return
		}
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	var parsedDate time.Time

	if err := utils.ParseDateToTime(input.BirthDate, &parsedDate); err != nil {
		c.JSON(http.StatusBadRequest, "incorrect type of birth date")
		return
	}

	user := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		BirthDate: parsedDate,
		Gender:    input.Gender,
		Email:     input.Email,
		Password:  string(hashedPassword),
	}

	database.DB.Create(&user)

	c.JSON(http.StatusCreated, user.ToSchema())
}

func LoginUser(c *gin.Context) {
	var input models.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RaiseBadRequestError(c, err)
		return
	}

	var user models.User

	if err := validators.ValidateUserLogin(database.DB, &user, input.Email, input.Password); err != nil {
		utils.RaiseBadRequestError(c, err)
		return
	}

	var token string

	if err := auth.GenerateToken(&token, user.ID); err != nil {
		utils.RaiseBadRequestError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"access": token})
}

func Profile(c *gin.Context) {
	userIDValue := c.Value("userID")

	if userIDValue == nil {
		c.JSON(http.StatusUnauthorized, "user ID is missing in context")
		return
	}

	userID, ok := userIDValue.(uint)

	if !ok {
		c.JSON(http.StatusUnauthorized, "user ID has wrong type")
		return
	}

	var user models.User
	database.DB.First(&user, userID)

	c.JSON(http.StatusOK, user.ToSchema())
}
