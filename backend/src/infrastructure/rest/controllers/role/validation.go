package role

import "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"

var customRules = map[string]string{
	"name":           "required",
	"default_router": "required",
	"order":          "required,numeric",
	"label":          "required",
	"description":    "required",
	"parent_id":      "required,lt=11",
	"status":         "required,status_enum",
}

func updateValidation(request map[string]any) error {
	validator := controllers.NewCommonValidator(customRules)
	return validator.ValidateUpdate(request)
}
