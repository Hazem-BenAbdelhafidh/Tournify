package tournament

type ITournamentService interface {
	CreateTournament(payload CreateTournament) (Tournament, error)
	UpdateTournament(id int, payload CreateTournament) error
	DeleteTournament(id int) error
	GetTournamentById(id int) (Tournament, error)
	GetTournaments(limit, offset int) ([]Tournament, error)
}

type TournamentService struct {
	repo ITournamentRepository
}

func NewTournamentService(repo ITournamentRepository) *TournamentService {
	return &TournamentService{
		repo: repo,
	}
}

func (ts *TournamentService) CreateTournament(payload CreateTournament) (Tournament, error) {
	return ts.repo.CreateTournament(payload)
}

func (ts *TournamentService) GetTournamentById(id int) (Tournament, error) {
	return ts.repo.GetTournamentById(id)
}

func (ts *TournamentService) GetTournaments(limit, offset int) ([]Tournament, error) {
	return ts.repo.GetTournaments(limit, offset)
}

func (ts *TournamentService) DeleteTournament(id int) error {
	return ts.repo.DeleteTournament(id)
}

func (ts *TournamentService) UpdateTournament(id int, payload CreateTournament) error {
	return ts.repo.DeleteTournament(id)
}
