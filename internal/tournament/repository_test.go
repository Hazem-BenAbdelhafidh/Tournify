package tournament_test

import (
	"testing"
	"time"

	"github.com/Hazem-BenAbdelhafidh/Tournify/db"
	"github.com/Hazem-BenAbdelhafidh/Tournify/entities"
	"github.com/Hazem-BenAbdelhafidh/Tournify/internal/tournament"
	"github.com/Hazem-BenAbdelhafidh/Tournify/utils"
	"github.com/stretchr/testify/suite"
)

type TournamentSuite struct {
	suite.Suite
	repo *tournament.TournamentRepo
}

func TestTournamentSuite(t *testing.T) {
	suite.Run(t, &TournamentSuite{})
}

func (ts *TournamentSuite) SetupSuite() {
	db := db.ConnectToDb()
	repo := tournament.NewTournamentRepo(db)
	ts.repo = repo
	err := ts.repo.DB.Exec("DELETE FROM tournaments").Error
	ts.Require().NoError(err)
}

func (ts *TournamentSuite) AfterTest(suiteTest, testName string) {
	clearTournamentsTable(ts)
}

func (ts *TournamentSuite) TestCreateTournament() {
	utils.CreateRandomTournament(ts.T(), ts.repo)
}

func (ts *TournamentSuite) TestGetTournamentById() {
	createdTournament := utils.CreateRandomTournament(ts.T(), ts.repo)

	tournament, err := ts.repo.GetTournamentById(int(createdTournament.ID))
	ts.Require().NoError(err)
	ts.Require().NotEmpty(tournament)
	ts.Require().Equal(createdTournament.ID, tournament.ID)
	ts.Require().Equal(createdTournament.Name, tournament.Name)
	ts.Require().Equal(createdTournament.Description, tournament.Description)
	ts.Require().Equal(createdTournament.Game, tournament.Game)
	ts.Require().Equal(createdTournament.NumOfTeams, tournament.NumOfTeams)
	ts.Require().WithinDuration(createdTournament.StartDate, tournament.StartDate, time.Second)
	ts.Require().WithinDuration(createdTournament.EndDate, tournament.EndDate, time.Second)

}

func (ts *TournamentSuite) TestGetTournaments() {
	createdTournaments := []entities.Tournament{}
	for i := 0; i < 10; i++ {
		createdTournament := utils.CreateRandomTournament(ts.T(), ts.repo)
		createdTournaments = append(createdTournaments, createdTournament)
	}

	tournaments, err := ts.repo.GetTournaments(10, 0)
	ts.Require().NoError(err)
	ts.Require().NotEmpty(tournaments)
	ts.Require().Len(tournaments, 10)
}

func (ts *TournamentSuite) TestDeleteTournament() {
	createdTournament := utils.CreateRandomTournament(ts.T(), ts.repo)
	err := ts.repo.DeleteTournament(int(createdTournament.ID))
	ts.Require().NoError(err)
	tournament, err := ts.repo.GetTournamentById(int(createdTournament.ID))
	ts.Require().NoError(err)
	ts.Require().Empty(tournament)
}

func (ts *TournamentSuite) TestUpdateTournament() {
	createdTournament := utils.CreateRandomTournament(ts.T(), ts.repo)
	updatePayload := tournament.CreateTournament{
		Name:        "hazem's tournament",
		Description: "new description test",
	}

	err := ts.repo.UpdateTournament(int(createdTournament.ID), updatePayload)
	ts.Require().NoError(err)
	tournament, err := ts.repo.GetTournamentById(int(createdTournament.ID))
	ts.Require().NoError(err)
	ts.Require().Equal(createdTournament.ID, tournament.ID)
	ts.Require().Equal(updatePayload.Name, tournament.Name)
	ts.Require().Equal(updatePayload.Description, tournament.Description)
	ts.Require().Equal(updatePayload.Game, tournament.Game)
	ts.Require().Equal(tournament.StartDate, tournament.StartDate)
	ts.Require().Equal(tournament.StartDate, tournament.StartDate)
	ts.Require().Equal(tournament.EndDate, tournament.EndDate)
}

func clearTournamentsTable(ts *TournamentSuite) {
	err := ts.repo.DB.Exec("DELETE FROM tournaments").Error
	ts.Require().NoError(err, "Error while cleaning Tournaments table")
}
