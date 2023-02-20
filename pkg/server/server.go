package server

import (
	"context"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/dgrijalva/jwt-go/v4/request"
	"github.com/gorilla/mux"
	"github.kostyadodich/demo/pkg/model"
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
	u.Use(jwtMiddleware)
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

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from request
		token, err := request.ParseFromRequest(
			r,
			request.AuthorizationHeaderExtractor,
			jwt.KnownKeyfunc(jwt.SigningMethodHS256, []byte("keykeykey")),
			request.WithClaims(&model.Claims{}),
		)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", token.Claims.(*model.Claims).ID)

		next.ServeHTTP(w, r.WithContext(ctx))
		return
	})
}
