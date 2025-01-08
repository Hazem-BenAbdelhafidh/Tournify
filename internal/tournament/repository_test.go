package tournament_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Hazem-BenAbdelhafidh/Tournify/db"
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
	fmt.Println("Here : ", testName)
	clearTournamentsTable(ts)
}

func (ts *TournamentSuite) CreateRandomTournament() tournament.Tournament {
	name := utils.RandomString(7)
	description := utils.RandomString(15)
	game := utils.RandomString(5)
	numOfTeams := utils.RandomNumber(10)
	tournament := tournament.CreateTournament{
		Name:        name,
		Description: description,
		NumOfTeams:  uint(numOfTeams),
		Game:        game,
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour * 72),
	}

	createdTournament, err := ts.repo.Create(tournament)
	ts.Require().NoError(err)
	ts.Require().NotEmpty(createdTournament)
	ts.Require().NotZero(createdTournament.ID)
	ts.Require().Equal(tournament.Name, createdTournament.Name)
	ts.Require().Equal(tournament.Description, createdTournament.Description)
	ts.Require().Equal(tournament.Game, createdTournament.Game)
	ts.Require().Equal(tournament.NumOfTeams, createdTournament.NumOfTeams)
	ts.Require().Equal(tournament.StartDate, createdTournament.StartDate)
	ts.Require().Equal(tournament.EndDate, createdTournament.EndDate)

	return createdTournament

}

func (ts *TournamentSuite) TestCreate() {
	ts.CreateRandomTournament()
}

func (ts *TournamentSuite) TestGetById() {
	createdTournament := ts.CreateRandomTournament()

	tournament, err := ts.repo.GetById(int(createdTournament.ID))
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

func (ts *TournamentSuite) TestGetAll() {
	createdTournaments := []tournament.Tournament{}
	for i := 0; i < 10; i++ {
		createdTournament := ts.CreateRandomTournament()
		createdTournaments = append(createdTournaments, createdTournament)
	}

	tournaments, err := ts.repo.GetAll(10, 0)
	ts.Require().NoError(err)
	ts.Require().NotEmpty(tournaments)
	ts.Require().Len(tournaments, 10)
}

func (ts *TournamentSuite) TestDelete() {
	createdTournament := ts.CreateRandomTournament()
	err := ts.repo.Delete(int(createdTournament.ID))
	ts.Require().NoError(err)
	tournament, err := ts.repo.GetById(int(createdTournament.ID))
	ts.Require().NoError(err)
	ts.Require().Empty(tournament)
}

func (ts *TournamentSuite) TestUpdate() {
	createdTournament := ts.CreateRandomTournament()
	updatePayload := tournament.CreateTournament{
		Name:        "hazem's tournament",
		Description: "new description test",
	}

	err := ts.repo.Update(int(createdTournament.ID), updatePayload)
	ts.Require().NoError(err)
	tournament, err := ts.repo.GetById(int(createdTournament.ID))
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
