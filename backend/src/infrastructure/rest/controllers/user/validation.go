package user

import "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"

var customRules = map[string]string{
	"user_name": "required,gt=3,lt=100",
	"email":     "required,email",
	"phone":     "required,custom_phone",
	"nick_name": "required",
	"status":    "required,status_enum",
}

func updateValidation(request map[string]any) error {
	validator := controllers.NewCommonValidator(customRules)
	return validator.ValidateUpdate(request)
}
