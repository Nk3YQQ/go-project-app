package models

import (
	"backend/internal/config"
)

type RegisterInput struct {
	FirstName       string              `json:"first_name"`
	LastName        string              `json:"last_name"`
	BirthDate       string              `json:"birth_date"`
	Gender          config.GenderChoice `json:"gender"`
	Email           string              `json:"email"`
	Password        string              `json:"password"`
	PasswordConfirm string              `json:"password_confirm"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSchema struct {
	ID        uint                `json:"id"`
	FirstName string              `json:"first_name"`
	LastName  string              `json:"last_name"`
	BirthDate string              `json:"birth_date"`
	Gender    config.GenderChoice `json:"gender"`
	Email     string              `json:"email"`
}

type ProjectCreateSchema struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	StartedAt   string              `json:"started_at"`
	Deadline    string              `json:"deadline"`
	Status      config.StatusChoice `json:"status"`
	Executors   []int               `json:"executors"`
	Tasks       []Task              `json:"tasks"`
}

type ProjectSchema struct {
	ID          uint                `json:"id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	StartedAt   string              `json:"started_at"`
	Deadline    string              `json:"deadline"`
	Status      config.StatusChoice `json:"status"`
	Executors   []UserSchema        `json:"executors"`
	Tasks       []Task              `json:"tasks"`
}

type ProjectUpdateSchema struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	StartedAt   string              `json:"started_at"`
	Deadline    string              `json:"deadline"`
	Status      config.StatusChoice `json:"status"`
	Executors   []int               `json:"executors"`
	Tasks       []Task              `json:"tasks"`
}

type TaskCreateSchema struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Deadline    string              `json:"deadline"`
	Status      config.StatusChoice `json:"status"`
	ProjectID   int                 `json:"project_id"`
	Executors   []int               `json:"executors"`
}

type TaskSchema struct {
	ID          uint                `json:"id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Deadline    string              `json:"deadline"`
	Status      config.StatusChoice `json:"status"`
	ProjectID   uint                `json:"project_id"`
	Executors   []UserSchema        `json:"executors"`
}

type TaskUpdateSchema struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Deadline    string              `json:"deadline"`
	Status      config.StatusChoice `json:"status"`
	ProjectID   int                 `json:"project_id"`
	Executors   []int               `json:"executors"`
}
