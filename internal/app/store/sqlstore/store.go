package sqlstore

import (
	"github.com/jmoiron/sqlx"
	"github.com/morozvol/AuthService/internal/app/store"
)

// Store ...
type Store struct {
	db             *sqlx.DB
	userRepository *UserRepository
}

// New ...
func New(dbPool *sqlx.DB) store.Store {
	return &Store{
		db: dbPool,
	}
}

// User ...
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
