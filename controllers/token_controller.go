package controllers

import (
	"errors"
	"jwtAuthGin/auth"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("supersecretkey")

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func Refresh(c *gin.Context) {

	tokenString := c.GetHeader("RefershToken")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "request does not contain an access token"})
		c.Abort()
		return
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(*JWTClaim)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("couldn't parse claims")})
		c.Abort()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("token expired")})
		c.Abort()
		return
	}

	accessTokenString, err := auth.GenerateAccessToken(claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	refreshTokenString, err := auth.GenerateRefreshToken(claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessTokenString,
		"refresh_token": refreshTokenString,
	})

}
