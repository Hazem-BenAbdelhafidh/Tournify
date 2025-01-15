package user

import (
	"time"

	"github.com/Hazem-BenAbdelhafidh/Tournify/internal/tournament"
)

type User struct {
	ID          uint                    `json:"id"`
	Username    string                  `json:"username" `
	Email       string                  `json:"email"`
	Password    string                  `json:"-" `
	Tournaments []tournament.Tournament `json:"tournaments" gorm:"foreignKey:CreatorId"`
	CreatedAt   time.Time               `json:"createdAt"`
	UpdatedAt   time.Time               `json:"updatedAt"`
	DeleteAt    *time.Time              `json:"deletedAt"`
}

type CreateUser struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUser struct {
	Username string `json:"username" `
	Email    string `json:"email"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignupResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
