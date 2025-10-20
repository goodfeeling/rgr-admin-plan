package menuParameter

import "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"

var customRules = map[string]string{
	"type":             "required,lt=191",
	"key":              "required,lt=191",
	"value":            "required,lt=191",
	"sys_base_menu_id": "required",
}

func updateValidation(request map[string]any) error {
	validator := controllers.NewCommonValidator(customRules)
	return validator.ValidateUpdate(request)
}
