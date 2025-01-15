package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Hazem-BenAbdelhafidh/Tournify/api"
	"github.com/Hazem-BenAbdelhafidh/Tournify/db"
	"github.com/Hazem-BenAbdelhafidh/Tournify/entities"
	"github.com/Hazem-BenAbdelhafidh/Tournify/internal/tournament"
	"github.com/Hazem-BenAbdelhafidh/Tournify/internal/user"
	"github.com/Hazem-BenAbdelhafidh/Tournify/utils"
	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite
	testingServer  *httptest.Server
	tournamentRepo *tournament.TournamentRepo
	userRepo       *user.UserRepo
	userService    user.IUserService
}

func (hs *HandlerSuite) SetupSuite() {
	router := api.SetupRouter()
	testingServer := httptest.NewServer(router)
	hs.testingServer = testingServer
	db := db.ConnectToDb()
	tournamentRepo := tournament.NewTournamentRepo(db)
	userRepo := user.NewUserRepo(db)
	hs.tournamentRepo = tournamentRepo
	hs.userRepo = userRepo
	hs.userService = user.NewUserService(userRepo)
	err := hs.tournamentRepo.DB.Exec("DELETE FROM tournaments").Error
	hs.Require().NoError(err)
	err = hs.userRepo.DB.Exec("DELETE FROM users").Error
	hs.Require().NoError(err)
}

func (hs *HandlerSuite) AfterTest(suiteTest, testName string) {
	clearTournamentsTable(hs)
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, &HandlerSuite{})
}

func (hs *HandlerSuite) TestCreatorMiddleware() {
	response1, err := signup([]byte(`{"username":"hazem","email":"hazem@gmail.com","password":"password123"}`), hs.testingServer.URL)
	defer response1.Body.Close()
	hs.Require().NoError(err)
	hs.Require().Equal(http.StatusCreated, response1.StatusCode)

	var tokenCookie string
	for _, cookie := range response1.Cookies() {
		if cookie.Name == "token" {
			tokenCookie = cookie.Value
		}
	}

	payload := tournament.CreateTournament{
		Name:        "tournament",
		Description: "description",
		NumOfTeams:  10,
		Game:        "game",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour * 24),
	}

	jsonBody, err := json.Marshal(payload)
	hs.Require().NoError(err)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/tournament", hs.testingServer.URL), bytes.NewReader(jsonBody))
	hs.Require().NoError(err)
	req.AddCookie(&http.Cookie{Name: "token", Value: tokenCookie, HttpOnly: true, Path: "/", Domain: "localhost", Secure: false})
	response2, err := hs.testingServer.Client().Do(req)
	defer response2.Body.Close()
	hs.Require().NoError(err)
	hs.Require().Equal(http.StatusCreated, response2.StatusCode)

	var createdTournament entities.Tournament
	err = json.NewDecoder(response2.Body).Decode(&createdTournament)
	hs.Require().NoError(err)
	hs.Require().Equal("tournament", createdTournament.Name)
	hs.Require().Equal("description", createdTournament.Description)
	hs.Require().Equal(uint(10), createdTournament.NumOfTeams)
	hs.Require().Equal("game", createdTournament.Game)

	req, err = http.NewRequest("DELETE", fmt.Sprintf("%s/tournament/%d", hs.testingServer.URL, createdTournament.ID), nil)
	hs.Require().NoError(err)
	req.AddCookie(&http.Cookie{Name: "token", Value: tokenCookie, HttpOnly: true, Path: "/", Domain: "localhost", Secure: false})
	response3, err := hs.testingServer.Client().Do(req)
	hs.Require().Equal(http.StatusOK, response3.StatusCode)

}

func (hs *HandlerSuite) TestGetTournamentById() {
	createdTournament := utils.CreateRandomTournament(hs.T(), hs.tournamentRepo)

	response, err := http.Get(fmt.Sprintf("%s/tournament/%d", hs.testingServer.URL, createdTournament.ID))
	defer response.Body.Close()
	hs.Require().NoError(err)
	hs.Require().Equal(http.StatusOK, response.StatusCode)

	var tournament entities.Tournament
	err = json.NewDecoder(response.Body).Decode(&tournament)
	hs.Require().NoError(err)
	hs.Require().EqualValues(createdTournament, tournament)
}

func (hs *HandlerSuite) TestCreateTournament() {
	jsonBody := []byte(`{"name": "tournament", "description": "description", "numOfTeams": 10, "game": "game", "startDate": "2021-01-01T00:00:00Z", "endDate": "2021-01-01T00:00:00Z"}`)
	bodyReader := bytes.NewReader(jsonBody)

	response, err := http.Post(fmt.Sprintf("%s/tournament/", hs.testingServer.URL), "application/json", bodyReader)
	defer response.Body.Close()
	hs.Require().NoError(err)
	hs.Require().Equal(http.StatusCreated, response.StatusCode)

	var tournament entities.Tournament
	err = json.NewDecoder(response.Body).Decode(&tournament)
	hs.Require().NoError(err)
	hs.Require().Equal("tournament", tournament.Name)
	hs.Require().Equal("description", tournament.Description)
	hs.Require().Equal(10, tournament.NumOfTeams)
	hs.Require().Equal("game", tournament.Game)
}

func clearTournamentsTable(hs *HandlerSuite) {
	err := hs.tournamentRepo.DB.Exec("DELETE FROM tournaments").Error
	hs.Require().NoError(err, "Error while cleaning Tournaments table")
	err = hs.userRepo.DB.Exec("DELETE FROM users").Error
	hs.Require().NoError(err, "Error while cleaning Users table")
}
