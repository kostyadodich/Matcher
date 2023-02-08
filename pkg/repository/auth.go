package repository

import (
	"database/sql"
	"github.kostyadodich/demo/pkg/model"
)

type AuthUser struct {
	db *sql.DB
}

func NewAuthUser(db *sql.DB) *AuthUser {
	return &AuthUser{db: db}
}

func (a *AuthUser) Create(user model.AuthUser) error {
	user.Password = user.GeneratePasswordHash(user.Password)

	_, err := a.db.Exec("INSERT INTO user_credentials (login, password) VALUES ($1, $2)",
		user.Login,
		user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthUser) CheckExist(login string, password string) (bool, error) {

	auth := model.AuthUser{}
	err := a.db.QueryRow(
		`SELECT login, password FROM user_credentials WHERE login=$1 AND password=$2`,
		login, password).Scan(
		&auth.Login, &auth.Password)
	if err == sql.ErrNoRows {
		return false, err
	} else if err != nil {
		return false, err
	}

	return true, nil
}
