package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const UnauthorizedMessage = "You don't have access"

func DecodeJwtToken(jwtToken string) (jwt.MapClaims, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		fmt.Println("Failed to parse token:", err)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func AuthMiddleware(c *gin.Context) {
	jwtCookie, err := c.Request.Cookie("token")
	if err != nil {
		respondWithError(c, http.StatusUnauthorized, errors.New(UnauthorizedMessage))
		return
	}

	jwtToken := jwtCookie.Value

	claims, err := DecodeJwtToken(jwtToken)
	if err != nil {
		respondWithError(c, http.StatusUnauthorized, errors.New(UnauthorizedMessage))
	}

	if userId, exists := claims["userId"]; exists {
		c.Set("userId", userId)
	} else {
		respondWithError(c, http.StatusUnauthorized, errors.New(UnauthorizedMessage))
	}

	c.Next()
}
