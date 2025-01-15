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

// GetTournamentById	godoc
// @Summary	Gets tournament by id
// @Description This endpoint is used to get a single tournament by id
// @Produce application/json
// @Success 200 {object} ResponseBody{}
// @Param id path int true "id of a tournament"
// @Tags tournament
// @Router /tournament/{id} [get]
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

	c.JSON(http.StatusOK, tournament)
}

// GetTournaments	godoc
// @Summary	Gets tournaments
// @Description This endpoint is used to get tournaments with pagination
// @Produce application/json
// @Success 200 {object} ResponseBody{}
// @Param limit query int false "Limit"
// @Param offset query int false "offset"
// @Tags tournament
// @Router /tournament [get]
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

	c.JSON(http.StatusOK, tournaments)
}

// CreaeteTournament	godoc
// @Summary	Create tournament
// @Description This endpoint is used to create a new tournament
// @Produce application/json
// @Success 201 {object} ResponseBody{}
// @Param CreatePayload body tournament.CreateTournament true "Tournament"
// @Tags tournament
// @Router /tournament [post]
func (th *TournamentHandler) CreateTournament(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	var tournamentToCreate tournament.CreateTournament
	err = c.ShouldBindJSON(&tournamentToCreate)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	tournamentToCreate.CreatorId = userId

	createdTournament, err := th.TournamentService.CreateTournament(tournamentToCreate)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, createdTournament)
}

// UpdateTournament	godoc
// @Summary	Update tournament
// @Description This endpoint is used to update a tournament
// @Produce application/json
// @Success 200 {object} ResponseBody{}
// @Param id path int true "id of a tournament"
// @Param UpdateTournament body tournament.CreateTournament true "Update Tournament Payload"
// @Tags tournament
// @Router /tournament/{id} [patch]
func (th *TournamentHandler) UpdateTournament(c *gin.Context) {
	var updatePayload tournament.CreateTournament
	id := c.Param("id")

	intId, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}
	err = c.ShouldBindJSON(&updatePayload)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	err = th.TournamentService.UpdateTournament(intId, updatePayload)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// DeleteTournament	godoc
// @Summary	Deletes tournament
// @Description This endpoint is used to delete a tournament
// @Produce application/json
// @Success 200 {object} ResponseBody{}
// @Param id path int true "id of a tournament"
// @Tags tournament
// @Router /tournament/{id} [delete]
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

	c.JSON(http.StatusOK, nil)
}
