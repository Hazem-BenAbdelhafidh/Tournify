package user_test

import (
	"strings"
	"testing"

	"github.com/Hazem-BenAbdelhafidh/Tournify/internal/user"
	"golang.org/x/crypto/bcrypt"
)

func (us *UserSuite) TestHashPassword() {
	password := "password"
	hashedPassword, err := user.HashPassword(password)
	us.Require().NoError(err)
	us.Require().NotEmpty(hashedPassword)
}

func (us *UserSuite) TestComparePassword() {
	password := "password"
	hashedPassword, err := user.HashPassword(password)
	us.Require().NoError(err)
	us.Require().NotEmpty(hashedPassword)

	err = user.ComparePassword(hashedPassword, password)
	us.Require().NoError(err)
}

func (us *UserSuite) TestLogin() {
	payload := user.CreateUser{
		Username: "username",
		Email:    "hazem@gmail.com",
		Password: "password",
	}

	resp, err := us.service.Signup(payload)
	us.Require().NoError(err)
	us.Require().NotEmpty(resp)
	us.Require().NotEmpty(resp.Token)
	us.Require().True(strings.HasPrefix(resp.Token, "eyJ"))

	tests := []struct {
		name    string
		payload user.LoginUser
		err     error
	}{
		{
			name: "valid login",
			payload: user.LoginUser{
				Email:    payload.Email,
				Password: payload.Password,
			},
			err: nil,
		},
		{
			name: "invalid login",
			payload: user.LoginUser{
				Email:    payload.Email,
				Password: "invalid",
			},
			err: bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, test := range tests {
		us.T().Run(test.name, func(t *testing.T) {
			token, err := us.service.Login(test.payload)
			us.Require().Equal(test.err, err)
			if err == nil {
				us.Require().NotEmpty(token)
				us.Require().True(strings.HasPrefix(token, "eyJ"))
			}

		})
	}
}

func (us *UserSuite) TestSignup() {
	payload := user.CreateUser{
		Username: "username",
		Email:    "hazem@gmail.com",
		Password: "password",
	}

	resp, err := us.service.Signup(payload)
	us.Require().NoError(err)
	us.Require().NotEmpty(resp)
	us.Require().NotEmpty(resp.Token)
	us.Require().True(strings.HasPrefix(resp.Token, "eyJ"))
}
