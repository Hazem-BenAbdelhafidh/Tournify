package user

import (
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(payload CreateUser) (User, error)
	UpdateUser(id int, payload CreateUser) error
	DeleteUser(id int) error
	GetUserById(id int) (User, error)
	GetUsers(limit, offset int, searchWord string) ([]User, error)
}

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(DB *gorm.DB) *UserRepo {
	return &UserRepo{
		DB: DB,
	}

}

func (ur UserRepo) GetUserById(id int) (User, error) {
	var user User

	err := ur.DB.First(&user, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return User{}, err
	}

	return user, nil
}

func (ur UserRepo) GetUsers(limit, offset int, searchWord string) ([]User, error) {
	var users []User

	query := ur.DB.Find(&users).Limit(limit).Offset(offset)
	if searchWord != "" {
		query.Where("email = ? OR username = ?", searchWord, searchWord)
	}

	err := query.Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return []User{}, err
	}

	return users, nil
}

func (ur UserRepo) CreateUser(payload CreateUser) (User, error) {
	user := User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}

	err := ur.DB.Create(&user).Error
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (ur UserRepo) UpdateUser(id int, payload CreateUser) error {
	user := User{
		ID:       uint(id),
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}

	err := ur.DB.Save(&user).Error
	if err != nil {
		return err
	}

	return nil

}

func (ur UserRepo) DeleteTournament(id int) error {
	err := ur.DB.Delete(&User{}, id).Error
	if err != nil {
		return err
	}

	return nil
}
