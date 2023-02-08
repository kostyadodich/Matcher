package handler

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go/v4"
	"github.kostyadodich/demo/pkg/model"
	"github.kostyadodich/demo/pkg/repository"
	"io"
	"net/http"
	"time"
)

type AuthUser struct {
	authRepo *repository.AuthUser
}

func NewAuthUser(authRepo *repository.AuthUser) *AuthUser {
	return &AuthUser{authRepo: authRepo}
}

func (a *AuthUser) SingIn(w http.ResponseWriter, r *http.Request) {
	user := model.AuthUser{}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(data, &user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.Password = user.GeneratePasswordHash(user.Password)

	row, err := a.authRepo.CheckExist(user.Login, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	} else if !row {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.At(time.Now())},
		Login: user.Login,
	})

	tokenRaw, err := token.SignedString([]byte("keykeykey"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(tokenRaw))
}

func (a *AuthUser) SingUp(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	auth := model.AuthUser{}
	if err = json.Unmarshal(data, &auth); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err = auth.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	if err = a.authRepo.Create(auth); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
