package postgres

import (
	"github.com/Prateek61/go_auth/graph/model"
	"github.com/go-pg/pg/v10"
)

type UsersRepo struct {
	DB *pg.DB
}

func (u *UsersRepo) GetUserByID(id string) (*model.User, error) {
	var user model.User
	err := u.DB.Model(&user).Where("id = ?", id).First()

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UsersRepo) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := u.DB.Model(&user).Where("username = ?", username).First()

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UsersRepo) GetUsers() ([]*model.User, error) {
	var users []*model.User
	err := u.DB.Model(&users).Select()

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UsersRepo) CreateUser(user *model.User) error {
	_, err := u.DB.Model(user).Insert()

	if err != nil {
		return err
	}

	return nil
}