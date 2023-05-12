package postgres

import (
	"github.com/Prateek61/go_auth/graph/model"
	"github.com/go-pg/pg/v10"
)

type UsersRepo struct {
	DB *pg.DB
}

func (u *UsersRepo) GetUserByField(field string, value string) (*model.User, error) {
	var user model.User
	err := u.DB.Model(&user).Where(field + " = ?", value).First()

	return &user, err
}

func (u *UsersRepo) GetUserByID(id string) (*model.User, error) {
	return u.GetUserByField("id", id)
}

func (u *UsersRepo) GetUserByUsername(username string) (*model.User, error) {
	return u.GetUserByField("username", username)
}

func (u *UsersRepo) GetUserByEmail(email string) (*model.User, error) {
	return u.GetUserByField("email", email)
}

func (u *UsersRepo) GetUsers() ([]*model.User, error) {
	var users []*model.User
	err := u.DB.Model(&users).Select()

	return users, err
}

func (u *UsersRepo) CreateUser(tx *pg.Tx, user *model.User) (*model.User, error) {
	_, err := tx.Model(user).Returning("*").Insert()

	return user, err
}