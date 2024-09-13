package models

import (
	"time"

	"backend/internal/config"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string
	LastName  string
	BirthDate time.Time
	Gender    config.GenderChoice
	Role      string
	Email     string `gorm:"unique"`
	Password  string
}

func (u *User) ToSchema() UserSchema {
	return UserSchema{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		BirthDate: u.BirthDate.Format("02.01.2006"),
		Gender:    u.Gender,
		Email:     u.Email,
	}
}

type Task struct {
	gorm.Model
	Title       string `gorm:"unique"`
	Description string
	Deadline    time.Time
	Status      config.StatusChoice `gorm:"default:created"`
	ProjectID   uint
	Executors   []User `gorm:"many2many:task_users"`
}

func (t *Task) ToSchema() TaskSchema {
	return TaskSchema{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Deadline:    t.Deadline.Format("01.06.2006"),
		Status:      t.Status,
		ProjectID:   t.ProjectID,
		Executors:   t.UsersToSchema(t.Executors),
	}
}

func (t *Task) UsersToSchema(users []User) []UserSchema {
	var serializedUsers []UserSchema

	for _, user := range users {
		serializedUsers = append(serializedUsers, user.ToSchema())
	}

	return serializedUsers
}

type Project struct {
	gorm.Model
	Title       string `gorm:"unique"`
	Description string
	StartedAt   time.Time
	Deadline    time.Time
	Status      config.StatusChoice `gorm:"default:created"`
	Executors   []User              `gorm:"many2many:project_users"`
	Tasks       []Task              `gorm:"many2many:project_tasks"`
}

func (p *Project) UsersToSchema(users []User) []UserSchema {
	var serializedUsers []UserSchema

	for _, user := range users {
		serializedUsers = append(serializedUsers, user.ToSchema())
	}

	return serializedUsers
}

func (p *Project) ToSchema() ProjectSchema {
	return ProjectSchema{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		StartedAt:   p.StartedAt.Format("01.06.2006"),
		Deadline:    p.StartedAt.Format("01.06.2006"),
		Status:      p.Status,
		Executors:   p.UsersToSchema(p.Executors),
		Tasks:       p.Tasks,
	}
}
