package sqlstore

import (
	"database/sql"
	"github.com/morozvol/AuthService/internal/app/model"
	"github.com/morozvol/AuthService/internal/app/store"
)

// UserRepository ...
type UserRepository struct {
	store *Store
}

// Create ...
func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRowx(
		"INSERT INTO users (email, encrypted_password, salt) VALUES ($1,$2, $3) RETURNING id",
		u.Email,
		u.EncryptedPassword,
		u.Salt,
	).Scan(&u.ID)
}

// Find ...
func (r *UserRepository) Find(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRowx(
		"SELECT id, email, encrypted_password FROM users WHERE id = $1",
		id,
	).StructScan(u); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

// FindByEmail ...
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}

	if err := r.store.db.QueryRowx(
		"SELECT id, email, encrypted_password, salt FROM users WHERE email = $1",
		email,
	).StructScan(u); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}
