package user

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Login(payload LoginUser) (string, error)
	Signup(payload CreateUser) (SignupResponse, error)
	DeleteUser(id int) error
	UpdateUser(id int, payload CreateUser) error
	GetUserById(id int) (User, error)
	GetUsers(limit, offset int, searchWord string) ([]User, error)
}

type UserService struct {
	repo IUserRepository
}

func NewUserService(repo IUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func signToken(userId uint) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
	}
	key := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return str, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil

}

func ComparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) Signup(payload CreateUser) (SignupResponse, error) {
	password := payload.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return SignupResponse{}, err
	}
	payload.Password = string(hashedPassword)
	user, err := us.repo.CreateUser(payload)
	if err != nil {
		return SignupResponse{}, err
	}

	token, err := signToken(user.ID)
	if err != nil {
		return SignupResponse{}, err
	}

	signupResponse := SignupResponse{
		Token: token,
		User:  user,
	}

	return signupResponse, nil
}

func (us *UserService) Login(payload LoginUser) (string, error) {
	user, err := us.repo.GetUserByEmail(payload.Email)
	if err != nil {
		return "", err
	}

	if user.ID == 0 {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return "", err
	}

	token, err := signToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (us *UserService) GetUserById(id int) (User, error) {
	return us.repo.GetUserById(id)
}

func (us *UserService) GetUsers(limit, offset int, searchWord string) ([]User, error) {
	return us.repo.GetUsers(limit, offset, searchWord)
}

func (us *UserService) DeleteUser(id int) error {
	return us.repo.DeleteUser(id)
}

func (us *UserService) UpdateUser(id int, payload CreateUser) error {
	return us.repo.UpdateUser(id, payload)
}
