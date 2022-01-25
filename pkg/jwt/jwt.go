// Package jwt  for issuing parsing and validating jwt tokens
package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/morozvol/AuthService/internal/app/model"
	"strconv"
	"time"
)

//MyCustomClaims  it is a payload jwt token
type MyCustomClaims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

//Valid test is valid payload
func (m *MyCustomClaims) Valid() error {
	return nil
	//TODO: вылидация
}

//GetClaims return pointer for base payload
func GetClaims(u *model.User) *MyCustomClaims {
	return &MyCustomClaims{
		"",
		jwt.StandardClaims{
			Id:        strconv.Itoa(u.ID),
			ExpiresAt: time.Now().Add(time.Hour * 100).Unix()},
	}
}

//NewJWT  create jwt token
func NewJWT(mySigningKey []byte, clams *MyCustomClaims) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims = clams
	tokenString, err := token.SignedString(mySigningKey)
	return tokenString, err
}

//ReadParse parse token
func ReadParse(myToken string, mySigningKey string) {
	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})

	if err == nil && token.Valid {
		fmt.Println("Your token is valid.  I like your style.")
	} else {
		fmt.Println("This token is terrible!  I cannot accept this.")
	}
}
