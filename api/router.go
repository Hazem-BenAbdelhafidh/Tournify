package api

import (
	"encoding/json"
	"io"

	"github.com/Hazem-BenAbdelhafidh/Tournify/db"
	"github.com/Hazem-BenAbdelhafidh/Tournify/internal/tournament"
	"github.com/Hazem-BenAbdelhafidh/Tournify/internal/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
)

func SetupRouter() *gin.Engine {

	db := db.ConnectToDb()
	// tournament
	tournamentRepo := tournament.NewTournamentRepo(db)
	tournamentService := tournament.NewTournamentService(tournamentRepo)
	tournamentHandler := NewTournamentHandler(tournamentService)

	// user
	userRepo := user.NewUserRepo(db)
	userService := user.NewUserService(userRepo)
	userHandler := NewUserHandler(userService)

	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("even", tournament.IsEven)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	tournamentRouter := router.Group("/tournament")
	{
		tournamentRouter.GET("", tournamentHandler.GetTournaments)
		tournamentRouter.GET("/:id", tournamentHandler.GetTournamentById)
		tournamentRouter.Use(AuthMiddleware)
		tournamentRouter.POST("/", tournamentHandler.CreateTournament)
		tournamentRouter.DELETE("/:id", tournamentHandler.DeleteTournament)
		tournamentRouter.PATCH("/:id", tournamentHandler.UpdateTournament)
	}

	userRouter := router.Group("/user")
	{
		userRouter.POST("/signup", userHandler.Signup)
		userRouter.POST("/login", userHandler.Login)
		userRouter.GET("")
		userRouter.GET("/:id")
		userRouter.Use(AuthMiddleware)
		userRouter.DELETE("/:id")
		userRouter.PATCH("/:id")
		userRouter.GET("/me", userHandler.GetMyInfo)
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
