package api

import (
	"encoding/json"
	"io"

	"github.com/Hazem-BenAbdelhafidh/Tournify/db"
	"github.com/Hazem-BenAbdelhafidh/Tournify/internal/tournament"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	db := db.ConnectToDb()
	tournamentRepo := tournament.NewTournamentRepo(db)
	tournamentService := tournament.NewTournamentService(tournamentRepo)
	tournamentHandler := NewTournamentHandler(tournamentService)
	router := gin.Default()

	tournamentRouter := router.Group("/tournament")
	{
		tournamentRouter.GET("", tournamentHandler.GetTournaments)
		tournamentRouter.GET("/:id", tournamentHandler.GetTournamentById)
		tournamentRouter.POST("/", tournamentHandler.CreateTournament)
		tournamentRouter.DELETE("/:id", tournamentHandler.DeleteTournament)
		tournamentRouter.PATCH("/:id", tournamentHandler.UpdateTournament)
	}

	userRouter := router.Group("/user")
	{
		userRouter.GET("")
		userRouter.GET("/:id")
		userRouter.POST("")
		userRouter.DELETE("/:id")
		userRouter.PATCH("/:id")
	}

	return router
}

type ResponseBody struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func respondWithJson(c *gin.Context, statusCode int, data any) {
	respBody := ResponseBody{
		Message: "success",
		Data:    data,
	}

	c.JSON(statusCode, respBody)
}

func (rb *ResponseBody) ToJson(w io.Writer) error {
	return json.NewEncoder(w).Encode(rb)
}

func (rb *ResponseBody) FromJson(r io.Reader) error {
	return json.NewDecoder(r).Decode(rb)
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
}

func (eb *ErrorResponse) ToJson(w io.Writer) error {
	return json.NewEncoder(w).Encode(eb)
}

func (eb *ErrorResponse) FromJson(r io.Reader) error {
	return json.NewDecoder(r).Decode(eb)
}

func respondWithError(c *gin.Context, statusCode int, err error) {
	errorResp := ErrorResponse{
		Message: "error",
		Error:   err,
	}

	c.JSON(statusCode, errorResp)
}
