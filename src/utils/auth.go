package utils

import (
	"context"
	"errors"
	models "java-gem/graph/model"
)

const UserContextKey string = "user"

func CheckUserRole(user *models.User, requiredRole models.UserRole) error {
	if user.Role != requiredRole {
		return errors.New("Forbidden access")
	}

	return nil
}

func GetUserIdFromContext(ctx context.Context) (string, error) {
	userId, ok := ctx.Value(UserContextKey).(string)
	if !ok || userId == "" {
		return "", errors.New("Unauthorized access")
	}

	return userId, nil

}
