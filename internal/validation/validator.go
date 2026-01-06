package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

var Validate *validator.Validate

func Init() {
	Validate = validator.New()
	pattern := `^([A-Za-z0-9\-._~]|%[0-9A-Fa-f]{2}|[!$&'()*+,;=]|[:@])*$` //rfc 9110 - path abempty -> rfc398 - segment
	_ = Validate.RegisterValidation("urlsegment", func(fl validator.FieldLevel) bool {
		token := fl.Field().String()
		matched, _ := regexp.MatchString(pattern, token)
		return matched
	})
	fmt.Println("Validator initialized")
}
