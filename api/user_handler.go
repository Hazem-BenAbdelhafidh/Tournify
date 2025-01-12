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

// GetUserById	godoc
// @Summary	gets a single user
// @Description This endpoint is used to get a single user by id
// @Param id path int true "User ID"
// @Produce application/json
// @Success 200 {object} ResponseBody{}
// @Tags user
// @Router /user/{id} [get]
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

// GetUsers	godoc
// @Summary	gets users
// @Description This endpoint is used to get users with pagination and search
// @Produce application/json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param search query string false "Search"
// @Success 200 {object} ResponseBody{}
// @Tags user
// @Router /user [get]
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

	users, err := uh.UserService.GetUsers(limit, offset, searchWord)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(c, http.StatusOK, users)
}

// Signup	godoc
// @Summary	signup
// @Description This endpoint is used to signup a new user using the username, email and password
// @Param CreatePayload body user.CreateUser true "User"
// @Produce application/json
// @Success 201 {object} ResponseBody{}
// @Tags user
// @Router /user/{id} [get]
func (uh *UserHandler) Signup(c *gin.Context) {
	var userToCreate user.CreateUser

	err := c.ShouldBindJSON(&userToCreate)
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

// Login	godoc
// @Summary	Login to an existent account
// @Description This endpoint is used to login to an existent account using the email and password
// @Produce application/json
// @Param LoginPayload body user.LoginUser true "User"
// @Success 200 {object} ResponseBody{}
// @Tags user
// @Router /user/login [post]
func (uh *UserHandler) Login(c *gin.Context) {
	var credentials user.LoginUser

	err := c.ShouldBindJSON(&credentials)
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

	respondWithJson(c, http.StatusOK, map[string]string{
		"token": token,
	})
}

// UpdateUser	godoc
// @Summary	update user
// @Description This endpoint is used to update existing user
// @Produce application/json
// @Success 200 {object} ResponseBody{}
// @Param id path int true "User ID"
// @Param UpdatePayload body user.UpdateUser true "User"
// @Tags user
// @Router /user/{id} [patch]
func (uh *UserHandler) UpdateUser(c *gin.Context) {
	var updatePayload user.UpdateUser
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

	err = uh.UserService.UpdateUser(intId, updatePayload)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(c, http.StatusOK, nil)
}

// DeleteUser	godoc
// @Summary	delete user
// @Description This endpoint is used to delete existing user
// @Produce application/json
// @Success 200 {object} ResponseBody{}
// @Param id path int true "User ID"
// @Tags user
// @Router /user/{id} [delete]
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

// GetMyInfo	godoc
// @Summary	gets the user info
// @Description This endpoint is used to get the current logged in user info
// @Produce application/json
// @Success 200 {object} ResponseBody{}
// @Tags user
// @Router /user/me [get]
func (uh *UserHandler) GetMyInfo(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		respondWithError(c, http.StatusUnauthorized, errors.New("No user id found"))
		return
	}

	floatUserId := userId.(float64)
	intUserId := int(floatUserId)

	user, err := uh.UserService.GetUserById(intUserId)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err)
		return
	}

	respondWithJson(c, http.StatusOK, user)
}
