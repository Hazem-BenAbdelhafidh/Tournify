package tournament

import (
	"gorm.io/gorm"
)

type TournamentRepository interface {
	CreateTournament(payload CreateTournament) (Tournament, error)
	UpdateTournament(id int, payload CreateTournament) error
	DeleteTournament(id int) error
	GetTournamentById(id int) (Tournament, error)
	GetTournaments(limit, offset int) ([]Tournament, error)
}

type TournamentRepo struct {
	DB *gorm.DB
}

func NewTournamentRepo(DB *gorm.DB) *TournamentRepo {
	return &TournamentRepo{
		DB: DB,
	}

}

func (tr TournamentRepo) GetTournamentById(id int) (Tournament, error) {
	var tournament Tournament

	err := tr.DB.First(&tournament, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return Tournament{}, err
	}

	return tournament, nil
}

func (tr TournamentRepo) GetTournaments(limit, offset int) ([]Tournament, error) {
	var tournaments []Tournament

	err := tr.DB.Find(&tournaments).Limit(limit).Offset(offset).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return []Tournament{}, err
	}

	return tournaments, nil
}

func (tr TournamentRepo) CreateTournament(payload CreateTournament) (Tournament, error) {
	tournament := Tournament{
		Name:        payload.Name,
		Description: payload.Description,
		NumOfTeams:  payload.NumOfTeams,
		Game:        payload.Game,
		StartDate:   payload.StartDate,
		EndDate:     payload.EndDate,
	}

	err := tr.DB.Create(&tournament).Error
	if err != nil {
		return Tournament{}, err
	}

	return tournament, nil
}

func (tr TournamentRepo) UpdateTournament(id int, payload CreateTournament) error {
	tournament := Tournament{
		ID:          uint(id),
		Name:        payload.Name,
		Description: payload.Description,
		NumOfTeams:  payload.NumOfTeams,
		Game:        payload.Game,
		StartDate:   payload.StartDate,
		EndDate:     payload.EndDate,
	}

	err := tr.DB.Save(&tournament).Error
	if err != nil {
		return err
	}

	return nil

}

func (tr TournamentRepo) DeleteTournament(id int) error {
	err := tr.DB.Delete(&Tournament{}, id).Error
	if err != nil {
		return err
	}

	return nil
}
