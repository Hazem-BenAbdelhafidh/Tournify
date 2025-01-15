package tournament

import "github.com/Hazem-BenAbdelhafidh/Tournify/entities"

type ITournamentService interface {
	CreateTournament(payload CreateTournament) (entities.Tournament, error)
	UpdateTournament(id int, payload CreateTournament) error
	DeleteTournament(id int) error
	GetTournamentById(id int) (entities.Tournament, error)
	GetTournaments(limit, offset int) ([]entities.Tournament, error)
}

type TournamentService struct {
	repo ITournamentRepository
}

func NewTournamentService(repo ITournamentRepository) *TournamentService {
	return &TournamentService{
		repo: repo,
	}
}

func (ts *TournamentService) CreateTournament(payload CreateTournament) (entities.Tournament, error) {
	return ts.repo.CreateTournament(payload)
}

func (ts *TournamentService) GetTournamentById(id int) (entities.Tournament, error) {
	return ts.repo.GetTournamentById(id)
}

func (ts *TournamentService) GetTournaments(limit, offset int) ([]entities.Tournament, error) {
	return ts.repo.GetTournaments(limit, offset)
}

func (ts *TournamentService) DeleteTournament(id int) error {
	return ts.repo.DeleteTournament(id)
}

func (ts *TournamentService) UpdateTournament(id int, payload CreateTournament) error {
	return ts.repo.DeleteTournament(id)
}
