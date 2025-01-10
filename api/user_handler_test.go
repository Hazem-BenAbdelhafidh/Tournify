package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Hazem-BenAbdelhafidh/Tournify/api"
	"github.com/Hazem-BenAbdelhafidh/Tournify/internal/user"
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

func (hs *HandlerSuite) Signup() *http.Response {
	jsonSignupBody := []byte(`{"username":"hazem","email":"hazem@gmail.com","password":"password123"}`)
	bodyReader := bytes.NewReader(jsonSignupBody)
	response, err := http.Post(fmt.Sprintf("%s/user/signup", hs.testingServer.URL), "application/json", bodyReader)
	hs.Require().NoError(err)
	hs.Require().Equal(http.StatusCreated, response.StatusCode)
	hs.Require().NotEmpty(response.Cookies())

	return response

}

func (hs *HandlerSuite) TestSignup() {
	response := hs.Signup()
	defer response.Body.Close()
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
	response := hs.Signup()
	defer response.Body.Close()
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
