package config

import (
	"github.com/go-playground/validator/v10"
	"log"
	"regexp"
	"strings"
)

func NewValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.RegisterValidation("not_contain_space", ValidateNotContainSpace); err != nil {
		log.Fatalln(err)
	}
	if err := v.RegisterValidation("is_password_format", ValidatePasswordFormat); err != nil {
		log.Fatalln(err)
	}
	return v
}

func ValidateNotContainSpace(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return !strings.Contains(field, " ")
}

func ValidatePasswordFormat(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",.<>\/?\\|~]`).MatchString(password)
	return hasLower && hasUpper && hasNumber && hasSpecial
}
