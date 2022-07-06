package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// func EnvLoad() {
// 	if err := godotenv.Load(); err != nil {
// 		log.Fatal(err)
// 	}
// }

// var jwtKey = os.Getenv("JWT_SECRET_KEY")

var jwtKey = []byte("supersecretkey")

type JWTClaim struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateAccessToken(email string) (accessTokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Minute)
	claims := &JWTClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err = token.SignedString(jwtKey)
	return
}

func GenerateRefreshToken(email string) (refreshTokenString string, err error) {
	expirationTime := time.Now().Add(2 * time.Minute)
	claims := &JWTClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return
}

func ValidateRefreshToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return
}
