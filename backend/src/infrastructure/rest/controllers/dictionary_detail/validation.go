package dictionary_detail

import "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"

var customRules = map[string]string{
	"label":  "required,lt=100",
	"type":   "required",
	"value":  "required",
	"extend": "omitempty",
	"status": "required,status_enum",
	"sort":   "omitempty,numeric",
}

func updateValidation(request map[string]any) error {
	validator := controllers.NewCommonValidator(customRules)
	return validator.ValidateUpdate(request)
}
