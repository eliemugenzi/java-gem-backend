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
	"java-gem/src/utils/constants"
	"java-gem/src/utils/validators"

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
	user := &models.User{}
	userResult := DB.Where("email = ?", email).First(user)
	if userResult.Error == nil {
		return nil, errors.New("User with this email already exists")
	}

	newUser := &models.User{
		ID:        uuid.New().String(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		CreatedAt: utils.GetCurrentTime(),
		UpdatedAt: utils.GetCurrentTime(),
		Password:  utils.HashPassword([]byte(password)),
		Role:      role,
	}

	DB.Create(newUser)

	tokenPair := utils.GenerateTokenPair(newUser.ID)
	newUser.Password = ""

	return &models.AuthPayload{
		Token: tokenPair["accessToken"],
		User:  newUser,
	}, nil

}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*models.AuthPayload, error) {

	foundUser := &models.User{}
	fmt.Println("User email: " + email)
	userResult := DB.Where("email = ?", email).First(foundUser)

	fmt.Println(userResult)

	if userResult.Error != nil {
		return nil, errors.New("Invalid credentials")
	}

	if !utils.ComparePassword(foundUser.Password, []byte(password)) {
		return nil, errors.New("Invalid credentials")
	}

	tokenPair := utils.GenerateTokenPair(foundUser.ID)

	foundUser.Password = ""

	return &models.AuthPayload{
		Token: tokenPair["accessToken"],
		User:  foundUser,
	}, nil
}

// CreateCoffee is the resolver for the createCoffee field.
func (r *mutationResolver) CreateCoffee(ctx context.Context, input *models.CreateCoffeeInput) (*models.Coffee, error) {
	err := validators.ValidateInput(input)
	if err != nil {
		fmt.Printf("Bad request: %s", err)
		return nil, err
	}
	userId := ctx.Value(constants.USER_CONTEXT_KEY)
	fmt.Printf("Logged in user ID: %s\n", userId)
	authenticatedUser := &models.User{}
	queryResult := DB.Where("id=?", userId).Take(authenticatedUser)

	if queryResult.Error != nil {
		errorMessage := fmt.Sprintf("Unable to find an authenticated user: %v", queryResult.Error.Error())
		return nil, errors.New(errorMessage)
	}

	if authenticatedUser.Role != models.UserRoleAdmin {
		errorMessage := fmt.Sprintf("You are not allowed to perform this operation")
		return nil, errors.New(errorMessage)
	}

	coffee := &models.Coffee{
		ID:          uuid.New().String(),
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		CreatedBy:   authenticatedUser,
		CreatedAt:   utils.GetCurrentTime(),
		UpdatedAt:   utils.GetCurrentTime(),
		Image:       input.Image,
	}

	DB.Create(coffee)
	coffee.CreatedBy.Password = ""

	return coffee, nil
}

func (r *queryResolver) CurrentUser(ctx context.Context) (*models.User, error) {
	userId := ctx.Value(constants.USER_CONTEXT_KEY)
	authenticatedUser := &models.User{}
	queryResult := DB.Where("id=?", userId).Take(authenticatedUser)

	if queryResult.Error != nil {
		errorMessage := fmt.Sprintf("Unable to find an authenticated user: %v", queryResult.Error.Error())
		return nil, errors.New(errorMessage)
	}

	authenticatedUser.Password = ""
	return authenticatedUser, nil
}

// CofeeList is the resolver for the cofeeList field.
func (r *queryResolver) CofeeList(ctx context.Context) ([]*models.Coffee, error) {
	coffeeList := []*models.Coffee{}
	queryResult := DB.Order("created_at desc").Preload("CreatedBy").Find(&coffeeList)

	if queryResult.Error != nil {
		errorMessage := fmt.Sprintf("Error when fetching coffee list: %v", queryResult.Error)
		return nil, errors.New(errorMessage)
	}

	return coffeeList, nil
}
