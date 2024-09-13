package validators

import (
	"backend/internal/models"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"regexp"
	"time"
)

type ErrorResponse struct {
	Details map[string]string `json:"details"`
}

func ValidateRegisterForm(firstName, lastName, email, password, PasswordConfirm string) map[string]string {
	validationErrors := make(map[string]string)

	if firstName == "" {
		validationErrors["first_name"] = "Поле 'first_name' обязательно"
	}

	if lastName == "" {
		validationErrors["last_name"] = "Поле 'last_name' обязательно"
	}

	if email == "" {
		validationErrors["email"] = "Поле 'email' обязательно"
	} else {
		if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
			validationErrors["email"] = "Формат поля 'email' некорректный"
		}
	}

	if len(password) < 6 {
		validationErrors["password"] = "Длина пароля не может быть меньше шести символов"
	}

	hasDigit := regexp.MustCompile(`[0-9]`).MatchString

	if !hasDigit(password) {
		validationErrors["password"] = "Пароль должен состоять из букв и цифр"
	}

	if password != PasswordConfirm {
		validationErrors["passwordConfirm"] = "Пароли не совпадают"
	}

	return validationErrors
}

func ValidateUserLogin(db *gorm.DB, user *models.User, email, password string) error {

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return errors.New("неверный логин или пароль")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return errors.New("неверный логин или пароль")
	}

	return nil
}

func ValidateDates(startedAt, deadline *time.Time) error {
	currentTime := time.Now()

	if deadline != nil && currentTime.After(*deadline) {
		return fmt.Errorf("current date cannot be later than the deadline")
	}

	if startedAt != nil && deadline != nil && deadline.Before(*startedAt) {
		return fmt.Errorf("deadline cannot be earlier than the start date")
	}

	return nil
}
