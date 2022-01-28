package token

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"time"
)

type Claims struct {
	UserName string
	jwt.StandardClaims

}

func CreateToken(Name ,SecretKey string) (string, error){

	KeyToken := []byte(SecretKey)

	claims := &Claims{
		UserName: Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(KeyToken)
	if err != nil {
		return "", errors.Wrap(err, "TOKEN SIGNING ERROR") //error handling
	}

	return tokenString, nil
	//json.NewEncoder(w).Encode(tokenString)
}

func CheckToken(HeaderAuth []string, jwtSecretKey string) (bool, error) {

	tokenString := HeaderAuth[1]

	claims := &Claims{}

	KeyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, KeyFunc)
	if err != nil {
		return false, errors.Wrap(err, "INCORRECT TOKEN SIGNATURE") //error handling
	}

	if !token.Valid {
		return false, nil
	}

	return true, nil
}
