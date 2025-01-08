package tournament

import "time"

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
	Name        string    `json:"name"`
	Description string    `json:"description"`
	NumOfTeams  uint      `json:"numOfTeams"`
	Game        string    `json:"game"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
}
