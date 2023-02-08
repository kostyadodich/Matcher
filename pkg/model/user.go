package model

import (
	"errors"
)

type User struct {
	ID           int64  `json:"id"`
	FullName     string `json:"full_name"`
	Age          int    `json:"age"`
	FavoriteDrug Item   `json:"favorite_drug"`
	Sex          bool   `json:"sex"`
}

func (u *User) IsMacancaAddicted() bool {
	return u.FavoriteDrug.ID == 1
}

func (u *User) IsMale() bool {
	return u.Sex == true
}

func (u *User) Validate() error {
	if u.Age < 0 || u.Age > 100 {
		return errors.New("incorrect age")
	}

	if u.IsMacancaAddicted() && u.Age > 30 {
		return errors.New("too old narc")
	}

	if len(u.FullName) < 3 {
		return errors.New("full name too short")
	}

	return nil
}
