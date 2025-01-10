package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Hazem-BenAbdelhafidh/Tournify/internal/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserService user.IUserService
}

func NewUserHandler(us *user.UserService) *UserHandler {
	return &UserHandler{
		UserService: us,
	}

}

func (uh *UserHandler) GetUserById(c *gin.Context) {
	id := c.Param("id")

	intId, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	user, err := uh.UserService.GetUserById(intId)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	respondWithJson(c, http.StatusOK, user)
}

func (uh *UserHandler) GetUsers(c *gin.Context) {
	var limit int
	var page int
	var searchWord string

	limit = c.GetInt("limit")
	if limit == 0 {
		limit = 10
	}

	page = c.GetInt("page")
	if page == 0 {
		page = 1
	}

	searchWord = c.Query("search")

	offset := (limit * page) - 1

	tournaments, err := uh.UserService.GetUsers(limit, offset, searchWord)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(c, http.StatusOK, tournaments)
}

func (uh *UserHandler) Signup(c *gin.Context) {
	var userToCreate user.CreateUser

	err := c.BindJSON(&userToCreate)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	createdUser, err := uh.UserService.Signup(userToCreate)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	c.SetCookie("token", createdUser.Token, 3600, "/", "localhost", false, true)

	respondWithJson(c, http.StatusCreated, map[string]string{
		"token": createdUser.Token,
	})
}

func (uh *UserHandler) Login(c *gin.Context) {
	var credentials user.LoginUser

	err := c.BindJSON(&credentials)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	token, err := uh.UserService.Login(credentials)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		respondWithError(c, http.StatusBadRequest, errors.New("Invalid credentials"))
		return
	} else if err != nil {
		respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	c.SetCookie("token", token, 3600, "/", "localhost", false, true)

	respondWithJson(c, http.StatusCreated, map[string]string{
		"token": token,
	})
}

func (uh *UserHandler) UpdateUser(c *gin.Context) {
	var updatePayload user.CreateUser
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

	err = uh.UserService.UpdateUser(intId, updatePayload)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(c, http.StatusOK, nil)
}

func (uh *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	intId, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	err = uh.UserService.DeleteUser(intId)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(c, http.StatusOK, nil)
}
