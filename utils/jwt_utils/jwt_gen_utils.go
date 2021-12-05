package jwt_utils

import (
	"os"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(claims jwt.MapClaims)(string, error){
	signKey := []byte(os.Getenv("SECRET_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtString, err := token.SignedString(signKey)
	if err != nil{
		return "", err
	}

	return jwtString, nil
}