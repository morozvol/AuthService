package apiserver

import (
	"encoding/json"
	"errors"
	"github.com/morozvol/AuthService/internal/app/model"
	"github.com/morozvol/AuthService/internal/app/store"
	"github.com/morozvol/AuthService/pkg/jwt"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
)

const (
	sessionName = "WebStore"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errUserWithThisEmailExists  = errors.New("a user with this email exists")
	errNotAuthenticated         = errors.New("not authenticated")
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
	config *Config
}

func newServer(store store.Store, config *Config) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
		config: config,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/login", s.handleLogin()).Methods("POST")
	s.router.HandleFunc("/signup", s.handleSignUp()).Methods("POST")
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) handleLogin() func(http.ResponseWriter, *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password, u.Salt) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}
		if token, err := jwt.NewJWT([]byte(s.config.JWTKey), jwt.GetClaims(u)); err == nil {
			s.respond(w, r, http.StatusOK, token)
			return
		}
		s.error(w, r, http.StatusInternalServerError, err)
	}
}

func (s *server) handleSignUp() func(http.ResponseWriter, *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if _, err := s.store.User().FindByEmail(req.Email); err == nil {
			s.error(w, r, http.StatusUnprocessableEntity, errUserWithThisEmailExists)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
		//TODO: при создании пользователя, при момощи RabbitMQ(например),
		// Отправляем всем сервисам для создания корины, профимя...
	}
}
