package utils

import (
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type jwtCustomClaim struct {
	UserId interface{} `json:"userId"`
	jwt.StandardClaims
}

func GetSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "fc93cb07e1ad92898527100e58a1cf1d1e7"
	}
	return secretKey
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKeyType
		}

		return []byte(GetSecretKey()), nil
	})
}

func GetTokenString(context *gin.Context) string {
	bearerToken := context.GetHeader("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}

func GetUserIdFromToken(token *jwt.Token) string {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims["userId"].(string)
	}

	return ""
}

func getTokenClaims(userId interface{}, expiryDays int) *jwtCustomClaim {
	return &jwtCustomClaim{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 0, expiryDays).Unix(),
			Issuer:    os.Getenv("TOKEN_ISSUER"),
			IssuedAt:  time.Now().Unix(),
		},
	}
}

func GenerateTokenPair(userId interface{}) map[string]string {
	tokenClaims := getTokenClaims(userId, 15)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	tokenString, err := token.SignedString([]byte(GetSecretKey()))

	if err != nil {
		panic(err)
	}

	refreshTokenClaims := getTokenClaims(userId, 30)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(GetSecretKey()))

	if err != nil {
		panic(err)
	}

	return map[string]string{
		"accessToken":  tokenString,
		"refreshToken": refreshTokenString,
	}
}
