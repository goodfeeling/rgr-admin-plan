package dictionary

import "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"

var customRules = map[string]string{
	"status":           "required,status_enum",
	"name":             "required,lt=100",
	"type":             "required,lt=100",
	"desc":             "omitempty,lt=200",
	"is_generate_file": "omitempty",
}

func updateValidation(request map[string]any) error {
	validator := controllers.NewCommonValidator(customRules)
	return validator.ValidateUpdate(request)
}
