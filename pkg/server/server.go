package server

import (
	"github.com/gorilla/mux"
	"github.kostyadodich/demo/pkg/repository"
	"github.kostyadodich/demo/pkg/server/handler"
	"github.kostyadodich/demo/pkg/service/matcher"
	"net/http"
)

type Server struct {
	userRepo       *repository.User
	dialerRepo     *repository.Dialer
	defaultMatcher *matcher.DefaultMatcher
	authUser       *repository.AuthUser
}

func NewServer(
	userRepo *repository.User,
	dialerRepo *repository.Dialer,
	defaultMatcher *matcher.DefaultMatcher,
	authUser *repository.AuthUser) *Server {

	return &Server{
		userRepo:       userRepo,
		dialerRepo:     dialerRepo,
		defaultMatcher: defaultMatcher,
		authUser:       authUser,
	}
}

func (s *Server) Serve() {
	var (
		userHandler     = handler.NewUser(s.userRepo)
		dialerHandler   = handler.NewDialer(s.dialerRepo)
		userAuthHandler = handler.NewAuthUser(s.authUser)
	)

	r := mux.NewRouter()

	u := r.PathPrefix("/user").Subrouter()
	u.HandleFunc("/", userHandler.GetUsers).Methods("GET")
	u.HandleFunc("/{id}", userHandler.GetUserByID).Methods("GET")
	u.HandleFunc("/", userHandler.CreateUser).Methods("POST")
	u.HandleFunc("/{id}", userHandler.UpdateUser).Methods("PUT")
	u.HandleFunc("/{id}", userHandler.DeleteUser).Methods("DELETE")

	d := r.PathPrefix("/dialer").Subrouter()
	d.HandleFunc("/", dialerHandler.GetDialers).Methods("GET")
	d.HandleFunc("/{id}", dialerHandler.GetDialerByID).Methods("GET")
	d.HandleFunc("/", dialerHandler.CreateDialer).Methods("POST")
	d.HandleFunc("/{id}", dialerHandler.UpdateDialer).Methods("PUT")
	d.HandleFunc("/{id}", dialerHandler.DeleteDialer).Methods("DELETE")

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/sign-up", userAuthHandler.SingUp).Methods("POST")
	auth.HandleFunc("/sign-in", userAuthHandler.SingIn).Methods("POST")

	http.ListenAndServe(":8080", r)
}
