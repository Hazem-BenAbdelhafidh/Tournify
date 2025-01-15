package user

import (
	"github.com/Hazem-BenAbdelhafidh/Tournify/entities"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(payload CreateUser) (entities.User, error)
	UpdateUser(id int, payload UpdateUser) error
	DeleteUser(id int) error
	GetUserById(id int) (entities.User, error)
	GetUserByEmail(email string) (entities.User, error)
	GetUsers(limit, offset int, searchWord string) ([]entities.User, error)
}

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(DB *gorm.DB) *UserRepo {
	return &UserRepo{
		DB: DB,
	}

}

func (ur UserRepo) GetUserById(id int) (entities.User, error) {
	var user entities.User

	err := ur.DB.First(&user, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return entities.User{}, err
	}

	return user, nil
}

func (ur UserRepo) GetUserByEmail(email string) (entities.User, error) {
	var user entities.User
	err := ur.DB.First(&user, "email = ?", email).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return entities.User{}, err
	}

	return user, nil

}

func (ur UserRepo) GetUsers(limit, offset int, searchWord string) ([]entities.User, error) {
	var users []entities.User

	query := ur.DB.Find(&users).Limit(limit).Offset(offset)
	if searchWord != "" {
		query.Where("email = ? OR username = ?", searchWord, searchWord)
	}

	err := query.Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return []entities.User{}, err
	}

	return users, nil
}

func (ur UserRepo) CreateUser(payload CreateUser) (entities.User, error) {
	user := entities.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}

	err := ur.DB.Create(&user).Error
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (ur UserRepo) UpdateUser(id int, payload UpdateUser) error {
	user := entities.User{
		ID:       uint(id),
		Username: payload.Username,
		Email:    payload.Email,
	}

	err := ur.DB.Save(&user).Error
	if err != nil {
		return err
	}

	return nil

}

func (ur UserRepo) DeleteUser(id int) error {
	err := ur.DB.Delete(&entities.User{}, id).Error
	if err != nil {
		return err
	}

	return nil
}
