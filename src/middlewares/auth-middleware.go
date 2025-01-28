package middlewares

import (
	"context"
	"fmt"
	"java-gem/src/utils"
	"java-gem/src/utils/constants"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func ValidateToken(token string) (*jwt.Token, error) {
	// Validate token
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(utils.GetSecretKey()), nil
	})

}

// AuthMiddleware verifies the JWT token and adds the user to the context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Extract the token
		bearerToken := strings.Split(authHeader, "Bearer ")
		if len(bearerToken) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		// Parse and validate the token
		token, err := ValidateToken(bearerToken[1])
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		// Extract claims
		userId := utils.GetUserIdFromToken(token)

		if userId == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		}

		// Add user to context
		ctx := context.WithValue(c.Request.Context(), constants.USER_CONTEXT_KEY, userId)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
