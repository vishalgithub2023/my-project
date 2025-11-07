package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CheckError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}

}

func CreateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"useremail": email,
			"exp":       time.Now().Add(time.Hour * 24).Unix(),
		})
	secretKey := []byte(os.Getenv("Secret_Key"))

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
