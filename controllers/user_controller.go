package controllers

import (
	"jwtAuthGin/auth"
	"jwtAuthGin/database"
	"jwtAuthGin/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterUser(context *gin.Context) {

	var user models.User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	_, err := database.ConnDB().Exec("INSERT INTO users (name,username,email,password) VALUES ($1,$2,$3,$4)", user.Name, user.Username, user.Email, user.Password)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, gin.H{"email": user.Email, "username": user.Username})
}

func Login(context *gin.Context) {

	var request TokenRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// check if email exists and password is correct
	row, err := database.ConnDB().Query("SELECT password FROM users WHERE email = $1", request.Email)
	if err != nil {
		context.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	var password string

	for row.Next() {
		if err := row.Scan(&password); err != nil {
			context.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	if password == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "beyle user yoook"})
		context.Abort()
		return
	}

	credentialError := models.CheckPassword(request.Password, password)
	if credentialError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		context.Abort()
		return
	}

	accessTokenString, err := auth.GenerateAccessToken(request.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	refreshTokenString, err := auth.GenerateRefreshToken(request.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"access_token":  accessTokenString,
		"refresh_token": refreshTokenString,
	})
}
