package entities

import (
	"time"
)

type User struct {
	ID            uint         `json:"id"`
	Username      string       `json:"username" `
	Email         string       `json:"email"`
	Password      string       `json:"-" `
	MyTournaments []Tournament `json:"myTournaments" gorm:"foreignKey:CreatorId"`
	Tournaments   []Tournament `json:"tournaments" gorm:"many2many:user_tournaments;"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`
	DeleteAt      *time.Time   `json:"deletedAt"`
}
