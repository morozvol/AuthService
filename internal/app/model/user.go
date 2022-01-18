package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User ...
type User struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
	Sold              string `json:"-"`
}

// Validate ...
func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(6, 100)),
	)
}

// BeforeCreate ...
func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		u.Sold = uuid.New().String()

		enc, err := encryptPassword(u.Password, u.Sold)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	}

	return nil
}

// Sanitize ...
func (u *User) Sanitize() {
	u.Password = ""
}

// ComparePassword ...
func (u *User) ComparePassword(password string, sold string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password+sold)) == nil
}

func encryptPassword(password string, sold string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password+sold), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
