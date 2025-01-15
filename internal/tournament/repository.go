package tournament

import (
	"github.com/Hazem-BenAbdelhafidh/Tournify/entities"
	"gorm.io/gorm"
)

type ITournamentRepository interface {
	CreateTournament(payload CreateTournament) (entities.Tournament, error)
	UpdateTournament(id int, payload CreateTournament) error
	DeleteTournament(id int) error
	GetTournamentById(id int) (entities.Tournament, error)
	GetTournaments(limit, offset int) ([]entities.Tournament, error)
}

type TournamentRepo struct {
	DB *gorm.DB
}

func NewTournamentRepo(DB *gorm.DB) *TournamentRepo {
	return &TournamentRepo{
		DB: DB,
	}

}

func (tr TournamentRepo) GetTournamentById(id int) (entities.Tournament, error) {
	var tournament entities.Tournament

	err := tr.DB.First(&tournament, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return entities.Tournament{}, err
	}

	return tournament, nil
}

func (tr TournamentRepo) GetTournaments(limit, offset int) ([]entities.Tournament, error) {
	var tournaments []entities.Tournament

	err := tr.DB.Find(&tournaments).Limit(limit).Offset(offset).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return []entities.Tournament{}, err
	}

	return tournaments, nil
}

func (tr TournamentRepo) CreateTournament(payload CreateTournament) (entities.Tournament, error) {
	tournament := entities.Tournament{
		Name:        payload.Name,
		Description: payload.Description,
		NumOfTeams:  payload.NumOfTeams,
		Game:        payload.Game,
		StartDate:   payload.StartDate,
		EndDate:     payload.EndDate,
		CreatorId:   uint(payload.CreatorId),
	}

	err := tr.DB.Create(&tournament).Error
	if err != nil {
		return entities.Tournament{}, err
	}

	return tournament, nil
}

func (tr TournamentRepo) UpdateTournament(id int, payload CreateTournament) error {
	tournament := entities.Tournament{
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
	err := tr.DB.Delete(&entities.Tournament{}, id).Error
	if err != nil {
		return err
	}

	return nil
}
