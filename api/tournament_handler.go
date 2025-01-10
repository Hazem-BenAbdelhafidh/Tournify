package api

import (
	"net/http"
	"strconv"

	"github.com/Hazem-BenAbdelhafidh/Tournify/internal/tournament"
	"github.com/gin-gonic/gin"
)

type TournamentHandler struct {
	TournamentService tournament.ITournamentService
}

func NewTournamentHandler(ts *tournament.TournamentService) *TournamentHandler {
	return &TournamentHandler{
		TournamentService: ts,
	}

}

func (th *TournamentHandler) GetTournamentById(c *gin.Context) {
	id := c.Param("id")

	intId, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	tournament, err := th.TournamentService.GetTournamentById(intId)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	respondWithJson(c, http.StatusOK, tournament)
}

func (th *TournamentHandler) GetTournaments(c *gin.Context) {
	var limit int
	var page int

	limit = c.GetInt("limit")
	if limit == 0 {
		limit = 10
	}

	page = c.GetInt("page")
	if page == 0 {
		page = 1
	}

	offset := (limit * page) - 1

	tournaments, err := th.TournamentService.GetTournaments(limit, offset)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(c, http.StatusOK, tournaments)
}

func (th *TournamentHandler) CreateTournament(c *gin.Context) {
	var tournamentToCreate tournament.CreateTournament

	err := c.BindJSON(&tournamentToCreate)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	createdTournament, err := th.TournamentService.CreateTournament(tournamentToCreate)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(c, http.StatusCreated, createdTournament)
}

func (th *TournamentHandler) UpdateTournament(c *gin.Context) {
	var updatePayload tournament.CreateTournament
	id := c.Param("id")

	intId, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}
	err = c.BindJSON(&updatePayload)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	err = th.TournamentService.UpdateTournament(intId, updatePayload)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(c, http.StatusOK, nil)
}

func (th *TournamentHandler) DeleteTournament(c *gin.Context) {
	id := c.Param("id")

	intId, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	err = th.TournamentService.DeleteTournament(intId)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(c, http.StatusOK, nil)
}
