package user_test

import (
	"testing"

	"github.com/Hazem-BenAbdelhafidh/Tournify/db"
	"github.com/Hazem-BenAbdelhafidh/Tournify/internal/user"
	"github.com/Hazem-BenAbdelhafidh/Tournify/utils"
	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
	repo    *user.UserRepo
	service *user.UserService
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, &UserSuite{})
}

func (us *UserSuite) SetupSuite() {
	db := db.ConnectToDb()
	repo := user.NewUserRepo(db)
	us.repo = repo
	us.service = user.NewUserService(repo)
	err := us.repo.DB.Exec("DELETE FROM users").Error
	us.Require().NoError(err)
}

func (us *UserSuite) AfterTest(suiteTest, testName string) {
	clearUsersTable(us)
}

func (us *UserSuite) TestCreateUser() {
	utils.CreateRandomUser(us.T(), us.repo)
}

func (us *UserSuite) TestGetUserById() {
	createdUser := utils.CreateRandomUser(us.T(), us.repo)

	user, err := us.repo.GetUserById(int(createdUser.ID))
	us.Require().NoError(err)
	us.Require().NotEmpty(user)
	us.Require().Equal(createdUser.ID, user.ID)
	us.Require().Equal(createdUser.Username, user.Username)
	us.Require().Equal(createdUser.Email, user.Email)
}

func (us *UserSuite) TestGetUserByEmail() {
	createdUser := utils.CreateRandomUser(us.T(), us.repo)

	user, err := us.repo.GetUserByEmail(createdUser.Email)
	us.Require().NoError(err)
	us.Require().NotEmpty(user)
	us.Require().Equal(createdUser.ID, user.ID)
	us.Require().Equal(createdUser.Username, user.Username)
	us.Require().Equal(createdUser.Email, user.Email)
}

func (us *UserSuite) TestGetUsers() {
	createdUsers := []user.User{}
	for i := 0; i < 10; i++ {
		createdUser := utils.CreateRandomUser(us.T(), us.repo)
		createdUsers = append(createdUsers, createdUser)
	}

	users, err := us.repo.GetUsers(10, 0, "")
	us.Require().NoError(err)
	us.Require().NotEmpty(users)
	us.Require().Len(users, 10)
}

func (us *UserSuite) TestDeleteUser() {
	createdUser := utils.CreateRandomUser(us.T(), us.repo)
	err := us.repo.DeleteUser(int(createdUser.ID))
	us.Require().NoError(err)
	tournament, err := us.repo.GetUserById(int(createdUser.ID))
	us.Require().NoError(err)
	us.Require().Empty(tournament)
}

func (us *UserSuite) TestUpdateUser() {
	createdUser := utils.CreateRandomUser(us.T(), us.repo)
	updatePayload := user.UpdateUser{
		Username: "new user name",
		Email:    "email2@yahoo.com",
	}

	err := us.repo.UpdateUser(int(createdUser.ID), updatePayload)
	us.Require().NoError(err)
	user, err := us.repo.GetUserById(int(createdUser.ID))
	us.Require().NoError(err)
	us.Require().Equal(createdUser.ID, user.ID)
	us.Require().Equal(updatePayload.Username, user.Username)
	us.Require().Equal(updatePayload.Email, user.Email)
}

func clearUsersTable(us *UserSuite) {
	err := us.repo.DB.Exec("DELETE FROM users").Error
	us.Require().NoError(err, "Error while cleaning users table")
}
