package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//IsAuth used to check user is authenticated or not
func IsAuth() gin.HandlerFunc {
	return checkToken(false)
}

//IsAdmin used to check user is admin role or not
func IsAdmin() gin.HandlerFunc {
	return checkToken(true)
}

func checkToken(adminOnly bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) < 2 {
			c.JSON(422, gin.H{
				"message": "Authorization token not provided",
			})

			c.Abort()
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("jwt_user_id", claims["user_id"])
			c.Set("jwt_is_agent", claims["user_role"])

			userRole := bool(claims["user_role"].(bool))

			if adminOnly && !userRole {
				c.JSON(403, gin.H{
					"status":  403,
					"message": "only agent allowed",
				})

				c.Abort()
				return
			}
		} else {
			c.JSON(422, gin.H{
				"message": "Invalid token",
				"error":   err,
			})

			c.Abort()
			return
		}
	}
}
