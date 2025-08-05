package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

var validate *validator.Validate

func Init() {
	validate = validator.New()
	pattern := `^([A-Za-z0-9\-._~]|%[0-9A-Fa-f]{2}|[!$&'()*+,;=]|[:@])*$` //rfc 9110 - path abempty -> rfc398 - segment
	_ = validate.RegisterValidation("urlsegment", func(fl validator.FieldLevel) bool {
		token := fl.Field().String()
		matched, _ := regexp.MatchString(pattern, token)
		return matched
	})
	fmt.Println("Validator initialized")
}
