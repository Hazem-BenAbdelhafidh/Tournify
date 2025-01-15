package user

import "github.com/Hazem-BenAbdelhafidh/Tournify/entities"

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
	User  entities.User `json:"user"`
	Token string        `json:"token"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
