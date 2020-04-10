package routes

import (
	"fmt"
	"os"
	"time"

	"github.com/bagus-aulia/custom-agent-allocation/models"
	"github.com/dgrijalva/jwt-go"
)

func createTokenCustomer(customer *models.Customer) string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   customer.ID,
		"user_role": false,
		"exp":       time.Now().AddDate(0, 0, 7).Unix(),
		"iat":       time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		fmt.Println(err)
	}

	return tokenString
}

func createTokenAgent(agent *models.Agent) string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   agent.ID,
		"user_role": true,
		"exp":       time.Now().AddDate(0, 0, 7).Unix(),
		"iat":       time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		fmt.Println(err)
	}

	return tokenString
}
