package tournament

import (
	"time"

	"github.com/go-playground/validator"
)

type Tournament struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	NumOfTeams  uint      `json:"numOfTeams"`
	Game        string    `json:"game"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   time.Time `json:"deletedAt"`
}

type CreateTournament struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	NumOfTeams  uint      `json:"numOfTeams" binding:"required,min=2,even"`
	Game        string    `json:"game" binding:"required"`
	StartDate   time.Time `json:"startDate" binding:"required"`
	EndDate     time.Time `json:"endDate" binding:"required"`
}

func IsEven(fl validator.FieldLevel) bool {
	input := fl.Field().Int()

	return input%2 == 0
}
