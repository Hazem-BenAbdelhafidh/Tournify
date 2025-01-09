package tournament

type TournamentService interface {
	CreateTournament(payload CreateTournament) (Tournament, error)
	UpdateTournament(id int, payload CreateTournament) error
	DeleteTournament(id int) error
	GetTournamentById(id int) (Tournament, error)
	GetTournaments(limit, offset int) ([]Tournament, error)
}

type TournamentServ struct {
	repo TournamentRepository
}

func NewTournamentService(repo TournamentRepository) *TournamentServ {
	return &TournamentServ{
		repo: repo,
	}
}

func (ts *TournamentServ) CreateTournament(payload CreateTournament) (Tournament, error) {
	return ts.repo.CreateTournament(payload)
}

func (ts *TournamentServ) GetTournamentById(id int) (Tournament, error) {
	return ts.repo.GetTournamentById(id)
}

func (ts *TournamentServ) GetTournaments(limit, offset int) ([]Tournament, error) {
	return ts.repo.GetTournaments(limit, offset)
}

func (ts *TournamentServ) DeleteTournament(id int) error {
	return ts.repo.DeleteTournament(id)
}

func (ts *TournamentServ) UpdateTournament(id int, payload CreateTournament) error {
	return ts.repo.DeleteTournament(id)
}
