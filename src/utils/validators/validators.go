package validators

import (
	"fmt"
	models "java-gem/graph/model"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationErrors []models.ValidationError

var validate *validator.Validate

func init() {
	validate = validator.New()

	validate.RegisterValidation("pasword_strengh", validatePasswordStrength)
}

/*
Format errors
*/
func (v ValidationErrors) Error() string {
	var errors []string
	for _, err := range v {
		errors = append(errors, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(errors, "; ")
}

func validatePasswordStrength(fieldLevel validator.FieldLevel) bool {
	password := fieldLevel.Field().String()

	// Check minimum requirements
	hasUpperCase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowerCase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)

	return hasUpperCase && hasLowerCase && hasNumber && hasSpecial
}

func ValidateInput(input interface{}) error {
	if err := validate.Struct(input); err != nil {
		var validationErrors ValidationErrors

		for _, err_ := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, models.ValidationError{
				Field:   err_.Field(),
				Message: getErrorMessage(err_),
			})
		}

		return validationErrors
	}
	return nil
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is requirefd"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Minimum length is %s", err.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s", err.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", err.Param())
	case "password_strength":
		return "Password must contain uppercase, lowercase, number, and special character"
	default:
		return fmt.Sprintf("Invalid value for %s", err.Field())
	}
}
