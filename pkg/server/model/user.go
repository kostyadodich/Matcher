package model

import "github.kostyadodich/demo/pkg/model"

type User struct {
	ID             int64  `json:"id"`
	FullName       string `json:"full_name"`
	Age            int    `json:"age"`
	FavoriteDrugID int64  `json:"favorite_drug_id"`
	Sex            bool   `json:"sex"`
}

func (u *User) Convert() model.User {
	return model.User{
		ID:           u.ID,
		FullName:     u.FullName,
		Age:          u.Age,
		FavoriteDrug: model.Item{ID: u.FavoriteDrugID},
		Sex:          u.Sex,
	}
}
