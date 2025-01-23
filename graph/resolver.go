package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import (
	"context"
	"errors"
	"fmt"
	"java-gem/config"
	models "java-gem/graph/model"
	"java-gem/src/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Resolver struct {
}

type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Password  string
}

var users = []models.User{}
var currentUser *User
var DB *gorm.DB = config.Configure()

//? Resolvers

// SignUp is the resolver for the signUp field.
func (r *mutationResolver) SignUp(ctx context.Context, firstName string, lastName string, email string, password string, role models.UserRole) (*models.AuthPayload, error) {

	for _, user := range users {
		if user.Email == email {
			return nil, errors.New("email already exists")
		}
	}

	fmt.Println("User role:  ", role)

	newUser := models.User{
		ID:        uuid.New().String(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		Password:  utils.HashPassword([]byte(password)),
		Role:      role,
	}

	DB.Create(&newUser)

	tokenPair := utils.GenerateTokenPair(newUser.ID)
	return &models.AuthPayload{
		Token: tokenPair["accessToken"],
		User:  &newUser,
	}, nil

}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*models.AuthPayload, error) {

	foundUser := models.User{}
	fmt.Println("User email: " + email)
	userResult := DB.Where("email = ?", email).First(&foundUser)

	fmt.Println(userResult)

	if userResult.Error != nil {
		return nil, errors.New("Invalid credentials")
	}

	if !utils.ComparePassword(foundUser.Password, []byte(password)) {
		return nil, errors.New("Invalid credentials")
	}

	tokenPair := utils.GenerateTokenPair(foundUser.ID)

	return &models.AuthPayload{
		Token: tokenPair["accessToken"],
		User:  &foundUser,
	}, nil
}

// CreateCoffee is the resolver for the createCoffee field.
func (r *mutationResolver) CreateCoffee(ctx context.Context, name string, description string, price float64) (*models.Coffee, error) {
	panic(fmt.Errorf("not implemented: CreateCoffee - createCoffee"))
}
