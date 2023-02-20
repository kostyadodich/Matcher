package handler

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.kostyadodich/demo/pkg/repository"
	req "github.kostyadodich/demo/pkg/server/model"
	"io"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	repo *repository.User
}

func NewUser(repo *repository.User) *User {
	return &User{
		repo: repo,
	}
}

func (u *User) GetUsers(w http.ResponseWriter, r *http.Request) {

	_ = r.Context().Value("user_id")

	context.Background()

	users, err := u.repo.List()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)

}

func (u *User) CreateUser(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userReq := req.User{}
	err = json.Unmarshal(data, &userReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	user := userReq.Convert()

	if err = user.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = u.repo.Create(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (u *User) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id"))
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	userReq := req.User{}
	err = json.Unmarshal(data, &userReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	user := userReq.Convert()

	if err = user.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	user.ID = int64(id)

	if err = u.repo.Update(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (u *User) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = u.repo.Delete(int64(id)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (u *User) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := u.repo.GetByID(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(data)
}
