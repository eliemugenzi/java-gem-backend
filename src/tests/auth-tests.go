package tests

import (
	"context"
	"java-gem/graph"
	models "java-gem/src/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockDB struct {
	users [] models.User
}

func TestSignUp(t *testing.T) {
	resolver := &graph.Resolver{ }


	t.Run("Successful signup", func (t *testing.T) {
		// Test for successful signup

		ctx := context.Background()
		firstName := "John"
		lastName := "Doe"
		email := "john.doe@example.com"
		password := "password123"

		authPayload, err := resolver.SignUp(ctx, firstName, lastName, email, password)
		assert.NoError(t, err)
		assert.NotNil(t, authPayload)
		assert.Equal(t, authPayload.User.FirstName, firstName)
	})
}
