package repository

import (
	"database/sql"
	"github.kostyadodich/demo/pkg/model"
)

type User struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *User {
	return &User{
		db: db,
	}
}

func (u *User) GetByID(id int64) (*model.User, error) {
	user := model.User{}
	err := u.db.QueryRow(
		`SELECT u.id, u.full_name, u.age, u.sex, i.id, i.name, i.is_legal 
		FROM users u JOIN items i ON u.favorite_drug_id = i.id
		WHERE u.id = $1`,
		id).
		Scan(&user.ID, &user.FullName, &user.Age, &user.Sex,
			&user.FavoriteDrug.ID, &user.FavoriteDrug.Name, &user.FavoriteDrug.IsLegal)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) List() ([]model.User, error) {
	rows, err := u.db.Query(
		`SELECT u.id, u.full_name, u.age, u.sex, i.id, i.name, i.is_legal
			 FROM users u JOIN items i 
			ON u.favorite_drug_id= i.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.User
	for rows.Next() {
		user := model.User{}
		err := rows.Scan(
			&user.ID, &user.FullName, &user.Age, &user.Sex,
			&user.FavoriteDrug.ID, &user.FavoriteDrug.Name, &user.FavoriteDrug.IsLegal)
		if err != nil {
			return nil, err
		}
		result = append(result, user)
	}

	return result, nil
}

func (u *User) Create(user model.User) error {

	_, err := u.db.Exec(
		"INSERT INTO users (full_name, age, sex, favorite_drug_id) VALUES ($1, $2, $3, $4)",
		user.FullName,
		user.Age,
		user.Sex,
		user.FavoriteDrug.ID,
	)

	return err
}

func (u *User) Update(user model.User) error {
	_, err := u.db.Exec("UPDATE users SET full_name= $1, age= $2, sex= $3, favorite_drug_id=$4 WHERE id = $5",
		user.FullName,
		user.Age,
		user.Sex,
		user.FavoriteDrug.ID,
		user.ID,
	)

	return err
}

func (u *User) Delete(id int64) error {
	_, err := u.db.Exec("DELETE FROM users WHERE id = $1", id)

	return err
}
