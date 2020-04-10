package routes

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bagus-aulia/custom-agent-allocation/config"
	"github.com/bagus-aulia/custom-agent-allocation/models"
	"github.com/danilopolani/gocialite/structs"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// RedirectHandler to correct oAuth URL
func RedirectHandler(c *gin.Context) {
	// Retrieve provider from route
	provider := c.Param("provider")

	providerSecrets := map[string]map[string]string{
		"github": {
			"clientID":     os.Getenv("CLIENT_ID_GITHUB"),
			"clientSecret": os.Getenv("CLIENT_SECRET_GITHUB"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/github/callback",
		},
		"google": {
			"clientID":     os.Getenv("CLIENT_ID_GOOGLE"),
			"clientSecret": os.Getenv("CLIENT_SECRET_GOOGLE"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/google/callback",
		},
	}

	providerScopes := map[string][]string{
		"github": []string{},
		"google": []string{},
	}

	providerData := providerSecrets[provider]
	actualScopes := providerScopes[provider]
	authURL, err := config.Gocial.New().
		Driver(provider).
		Scopes(actualScopes).
		Redirect(
			providerData["clientID"],
			providerData["clientSecret"],
			providerData["redirectURL"],
		)

	// Check for errors (usually driver not valid)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	// Redirect with authURL
	c.Redirect(http.StatusFound, authURL)
}

// CallbackHandler callback of provider
func CallbackHandler(c *gin.Context) {
	// Retrieve query params for state and code
	state := c.Query("state")
	code := c.Query("code")
	provider := c.Param("provider")

	// Handle callback and check for errors
	user, _, err := config.Gocial.Handle(state, code)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	newAgent := getOrRegisterAgent(provider, user)
	jwtToken := createTokenAgent(&newAgent)

	c.JSON(200, gin.H{
		"data":    newAgent,
		"token":   jwtToken,
		"user":    user,
		"message": "login success",
	})
}

func getOrRegisterAgent(provider string, agent *structs.User) models.Agent {
	var userData models.Agent
	username := agent.Username

	config.DB.Where("provider = ? AND social_id = ?", provider, agent.ID).First(&userData)

	if username == "" {
		username = agent.ID
	}

	if userData.ID == 0 {
		newAgent := models.Agent{
			Nama:     username,
			Email:    agent.Email,
			SocialID: agent.ID,
			Provider: provider,
			Avatar:   agent.Avatar,
		}

		config.DB.Create(&newAgent)
		return newAgent
	}

	return userData
}

func createTokenAgent(agent *models.Agent) string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   agent.ID,
		"user_role": 1,
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
