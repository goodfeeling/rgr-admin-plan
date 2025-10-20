package api

import "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"

var customRules = map[string]string{
	"path":        "required,gt=2,lt=255",
	"method":      "required",
	"api_group":   "required",
	"description": "omitempty",
}

func updateValidation(request map[string]any) error {
	validator := controllers.NewCommonValidator(customRules)
	return validator.ValidateUpdate(request)
}
