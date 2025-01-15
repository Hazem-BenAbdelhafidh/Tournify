package utils

import (
	"math/rand/v2"
	"strings"
	"testing"
	"time"

	"github.com/Hazem-BenAbdelhafidh/Tournify/entities"
	"github.com/Hazem-BenAbdelhafidh/Tournify/internal/tournament"
	"github.com/Hazem-BenAbdelhafidh/Tournify/internal/user"
	"github.com/stretchr/testify/require"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomString(strLen int) string {
	var sb strings.Builder
	alphabetLength := len(alphabet)

	for i := 0; i < strLen; i++ {
		randomNumber := rand.IntN(alphabetLength)
		sb.WriteByte(alphabet[randomNumber])
	}

	return sb.String()
}

func RandomNumber(upperBound int) int {
	randomNumber := rand.IntN(upperBound)
	return randomNumber
}

func CreateRandomTournament(t *testing.T, repo *tournament.TournamentRepo) entities.Tournament {
	name := RandomString(7)
	description := RandomString(15)
	game := RandomString(5)
	numOfTeams := RandomNumber(10)
	tournament := tournament.CreateTournament{
		Name:        name,
		Description: description,
		NumOfTeams:  uint(numOfTeams),
		Game:        game,
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour * 72),
	}

	createdTournament, err := repo.CreateTournament(tournament)
	require.NoError(t, err)
	require.NotEmpty(t, createdTournament)
	require.NotZero(t, createdTournament.ID)
	require.Equal(t, tournament.Name, createdTournament.Name)
	require.Equal(t, tournament.Description, createdTournament.Description)
	require.Equal(t, tournament.Game, createdTournament.Game)
	require.Equal(t, tournament.NumOfTeams, createdTournament.NumOfTeams)
	require.Equal(t, tournament.StartDate, createdTournament.StartDate)
	require.Equal(t, tournament.EndDate, createdTournament.EndDate)

	return createdTournament

}

func CreateRandomUser(t *testing.T, repo *user.UserRepo) entities.User {
	username := RandomString(7)
	email := RandomString(15)
	password := RandomString(9)
	user := user.CreateUser{
		Username: username,
		Email:    email,
		Password: password,
	}

	createdUser, err := repo.CreateUser(user)
	require.NoError(t, err)
	require.NotEmpty(t, createdUser)
	require.NotZero(t, createdUser.ID)
	require.Equal(t, user.Username, createdUser.Username)
	require.Equal(t, user.Email, createdUser.Email)

	return createdUser
}
