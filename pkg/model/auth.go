package model

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
)

const salt = "0228papirosim"

type AuthUser struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Claims struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
	jwt.StandardClaims
}

func (a *AuthUser) Validate() error {
	if len(a.Password) < 4 {
		return errors.New("password is short")
	}

	return nil
}

func (a *AuthUser) GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password + salt))

	return fmt.Sprintf("%x", hash.Sum([]byte(password)))
}
