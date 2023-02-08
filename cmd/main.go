package main

import (
	"database/sql"
	"fmt"
	"github.kostyadodich/demo/pkg/repository"
	"github.kostyadodich/demo/pkg/server"
	"github.kostyadodich/demo/pkg/service/matcher"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d database=%s sslmode=disable",
		"postgres",
		"1234",
		"localhost",
		5432,
		"narcos",
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUser(db)
	dialerRepo := repository.NewDialer(db)
	DefaultMatcher := matcher.NewDefaultMatcher(db)
	authUser := repository.NewAuthUser(db)

	server.NewServer(userRepo, dialerRepo, DefaultMatcher, authUser).Serve()

}
