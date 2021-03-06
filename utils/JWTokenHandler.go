package utils

import (
	"fmt"
	"net/http"
	"projects/entitys"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var secretKey string = "Y29uZHVvbmdzdWFlbWRpODkzNA=="

type JwtCustomClaims struct {
	Id   int       `json:"id"`
	Name string    `json:"name"`
	Uuid uuid.UUID `json:"uuid"`
	jwt.StandardClaims
}

func GenerateRefreshToken(Id int, userName string) entitys.RefreshToken {
	// Set custom claims
	var refreshToken entitys.RefreshToken
	refreshToken.UserName = userName
	uuid := uuid.New()
	expiresAt := time.Now().Add(time.Hour * 12).Unix()
	refreshToken.ExpiresAt = expiresAt
	claims := &JwtCustomClaims{
		Id,
		userName,
		uuid,
		jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	result, err := token.SignedString([]byte(secretKey))
	refreshToken.Token = result
	if err != nil {
		return entitys.RefreshToken{}
	} else {
		return refreshToken
	}

}
func GenerateJWT(Id int, userName string) string {
	uuid := uuid.New()
	// Set custom claims
	claims := &JwtCustomClaims{
		Id,
		userName,
		uuid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		},
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	result, _ := token.SignedString([]byte(secretKey))

	return result

}

func ParseToken(c echo.Context) error {
	var userName string
	tokenString := c.Param("token")
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		userName = claims.Name
	} else {
		fmt.Println(err)
	}
	return c.String(http.StatusOK, userName)
}

func GetUserFromTokden(tokenString string) entitys.User {
	var user entitys.User
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		user.Name = claims.Name
		user.Id = claims.Id
	} else {
		//return exceptions.InValidTokenException(c)
		fmt.Println(err)
	}
	return user
}

func IsValidToken(validToken string) bool {

	token, er := jwt.ParseWithClaims(validToken, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if er != nil {
		return false
	}
	return token.Valid
}

func IsExpireToken(tokenString string) bool {
	var expireTime int64
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		expireTime = claims.StandardClaims.ExpiresAt
	} else {
		fmt.Println(err)
	}
	now := time.Now().Unix()
	if now < expireTime {
		return true
	} else {
		return false
	}
}
