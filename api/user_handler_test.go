package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Hazem-BenAbdelhafidh/Tournify/api"
	"github.com/Hazem-BenAbdelhafidh/Tournify/internal/user"
	"github.com/Hazem-BenAbdelhafidh/Tournify/utils"
)

const TestEmail = "hazem@gmail.com"

func (hs *HandlerSuite) TestDecodeJwtToken() {
	signupPayload := user.CreateUser{
		Username: "hazem",
		Email:    TestEmail,
		Password: "password123",
	}
	userResp, err := hs.userService.Signup(signupPayload)
	hs.Require().NoError(err)
	hs.Require().NotEmpty(userResp)

	claims, err := api.DecodeJwtToken(userResp.Token)
	hs.Require().NoError(err)
	hs.Require().Contains(claims, "userId")
}

func signup(payload []byte, url string) (*http.Response, error) {
	bodyReader := bytes.NewReader(payload)
	return http.Post(fmt.Sprintf("%s/user/signup", url), "application/json", bodyReader)
}

func (hs *HandlerSuite) TestSignup() {
	tests := []struct {
		name    string
		payload []byte
		error   bool
	}{
		{
			name:    "valid payload",
			payload: []byte(`{"username":"hazem","email":"hazem@gmail.com","password":"password123"}`),
			error:   false,
		},
		{
			name:    "missing username",
			payload: []byte(`{"email":"hazem@gmail.com","password":"password123"}`),
			error:   true,
		},
		{
			name:    "invalid email",
			payload: []byte(`{"username":"hazem","email":"hazemgmail.com","password":"password123"}`),
			error:   true,
		},
	}

	for _, test := range tests {
		response, err := signup(test.payload, hs.testingServer.URL)
		defer response.Body.Close()
		if test.error {
			hs.Require().Equal(http.StatusBadRequest, response.StatusCode)
		} else {
			hs.Require().NoError(err)
			hs.Require().NotNil(response)
			hs.Require().Equal(http.StatusCreated, response.StatusCode)
			hs.Require().NotEmpty(response.Cookies())
			tokenCookieExists := false
			for _, cookie := range response.Cookies() {
				if cookie.Name == "token" {
					tokenCookieExists = true
				}
			}
			hs.Require().True(tokenCookieExists)

			rb := api.ResponseBody{}
			err := rb.FromJson(response.Body)
			hs.Require().NoError(err)
			hs.Require().Equal("success", rb.Message)
			data, ok := rb.Data.(map[string]interface{})
			hs.Require().True(ok)
			hs.Require().NotEmpty(data["token"])
		}
	}
}

func (hs *HandlerSuite) TestLogin() {
	signupPayload := user.CreateUser{
		Username: "hazem",
		Email:    TestEmail,
		Password: "password123",
	}
	userResp, err := hs.userService.Signup(signupPayload)
	hs.Require().NoError(err)
	hs.Require().NotEmpty(userResp)

	jsonSignupBody := []byte(`{"email":"hazem@gmail.com","password":"password123"}`)
	bodyReader := bytes.NewReader(jsonSignupBody)
	response, err := http.Post(fmt.Sprintf("%s/user/login", hs.testingServer.URL), "application/json", bodyReader)
	defer response.Body.Close()
	hs.Require().NoError(err)
	hs.Require().Equal(http.StatusOK, response.StatusCode)
	hs.Require().NotEmpty(response.Cookies())
	tokenCookieExists := false
	for _, cookie := range response.Cookies() {
		if cookie.Name == "token" {
			tokenCookieExists = true
		}
	}
	hs.Require().True(tokenCookieExists)

	rb := api.ResponseBody{}
	err = rb.FromJson(response.Body)
	hs.Require().NoError(err)
	hs.Require().Equal("success", rb.Message)
	data, ok := rb.Data.(map[string]interface{})
	hs.Require().True(ok)
	hs.Require().NotEmpty(data["token"])
}

func (hs *HandlerSuite) TestMyInfo() {
	response, err := signup([]byte(`{"username":"hazem","email":"hazem@gmail.com","password":"password123"}`), hs.testingServer.URL)
	defer response.Body.Close()
	hs.Require().NoError(err)
	var tokenCookie string

	for _, cookie := range response.Cookies() {
		if cookie.Name == "token" {
			tokenCookie = cookie.Value
		}
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user/me", hs.testingServer.URL), nil)
	hs.Require().NoError(err)
	req.AddCookie(&http.Cookie{Name: "token", Value: tokenCookie, HttpOnly: true, Path: "/", Domain: "localhost", Secure: false})
	response, err = hs.testingServer.Client().Do(req)
	defer response.Body.Close()
	hs.Require().NoError(err)
	hs.Require().Equal(http.StatusOK, response.StatusCode)

	rb := api.ResponseBody{}
	err = rb.FromJson(response.Body)
	hs.Require().NoError(err)
	hs.Require().Equal("success", rb.Message)
	var user user.User
	data, err := json.Marshal(rb.Data)
	hs.Require().NoError(err)
	hs.Require().NotEmpty(data)
	err = json.Unmarshal(data, &user)
	hs.Require().NoError(err)

	hs.Require().Equal("hazem", user.Username)
	hs.Require().Equal(TestEmail, user.Email)
}

func (hs *HandlerSuite) TestGetUsers() {
	for i := 0; i < 10; i++ {
		utils.CreateRandomUser(hs.T(), hs.userRepo)
	}

	response, err := http.Get(fmt.Sprintf("%s/user", hs.testingServer.URL))
	defer response.Body.Close()
	hs.Require().NoError(err)
	hs.Require().Equal(http.StatusOK, response.StatusCode)

	rb := api.ResponseBody{}
	err = rb.FromJson(response.Body)
	hs.Require().NoError(err)
	hs.Require().Equal("success", rb.Message)
	data := rb.Data.([]interface{})
	hs.Require().Len(data, 10)
}

func (hs *HandlerSuite) TestGetUserById() {
	createdUser := utils.CreateRandomUser(hs.T(), hs.userRepo)

	response, err := http.Get(fmt.Sprintf("%s/user/%d", hs.testingServer.URL, createdUser.ID))
	defer response.Body.Close()
	hs.Require().NoError(err)
	hs.Require().Equal(http.StatusOK, response.StatusCode)

	rb := api.ResponseBody{}
	err = rb.FromJson(response.Body)
	hs.Require().NoError(err)
	hs.Require().Equal("success", rb.Message)
	data, err := json.Marshal(rb.Data)
	hs.Require().NoError(err)

	var user user.User
	err = json.Unmarshal(data, &user)
	hs.Require().NoError(err)
	hs.Require().Equal(createdUser.ID, user.ID)
	hs.Require().Equal(createdUser.Username, user.Username)
	hs.Require().Equal(createdUser.Email, user.Email)
	hs.Require().Equal("", user.Password)
}

func (hs *HandlerSuite) TestDeleteUser() {
	response, err := signup([]byte(`{"username": "hazem2","email":"hazem@gmail.com","password":"password123"}`), hs.testingServer.URL)
	defer response.Body.Close()

	var tokenCookie string
	for _, cookie := range response.Cookies() {
		if cookie.Name == "token" {
			tokenCookie = cookie.Value
		}
	}

	createdUser := utils.CreateRandomUser(hs.T(), hs.userRepo)
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/user/%d", hs.testingServer.URL, createdUser.ID), nil)
	hs.Require().NoError(err)
	req.AddCookie(&http.Cookie{Name: "token", Value: tokenCookie, HttpOnly: true, Path: "/", Domain: "localhost", Secure: false})
	response, err = hs.testingServer.Client().Do(req)
	defer response.Body.Close()
	hs.Require().NoError(err)
	hs.Require().Equal(http.StatusOK, response.StatusCode)
}

func (hs *HandlerSuite) TestUpdateUser() {
	response, err := signup([]byte(`{"username": "hazem2","email":"hazem@gmail.com","password":"password123"}`), hs.testingServer.URL)
	defer response.Body.Close()
	hs.Require().NoError(err)

	var tokenCookie string
	for _, cookie := range response.Cookies() {
		if cookie.Name == "token" {
			tokenCookie = cookie.Value
		}
	}

	createdUser := utils.CreateRandomUser(hs.T(), hs.userRepo)

	updatePayload := user.UpdateUser{
		Username: "hazem69",
		Email:    "hazem222@gmail.com",
	}

	jsonUpdatePayload, err := json.Marshal(updatePayload)
	hs.Require().NoError(err)

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/user/%d", hs.testingServer.URL, createdUser.ID), bytes.NewReader(jsonUpdatePayload))
	hs.Require().NoError(err)
	req.AddCookie(&http.Cookie{Name: "token", Value: tokenCookie, HttpOnly: true, Path: "/", Domain: "localhost", Secure: false})
	response, err = hs.testingServer.Client().Do(req)
	defer response.Body.Close()
	hs.Require().NoError(err)
	hs.Require().Equal(http.StatusOK, response.StatusCode)
}
