package menu

import "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"

var customRules = map[string]string{
	"component":  "required,lt=191",
	"title":      "required,lt=191",
	"name":       "required,lt=191",
	"path":       "required,lt=191",
	"hidden":     "omitempty",
	"keep_alive": "omitempty",
	"parent_id":  "omitempty,min=0",
	"icon":       "required,lt=191",
	"sort":       "omitempty",
}

func updateValidation(request map[string]any) error {
	validator := controllers.NewCommonValidator(customRules)
	return validator.ValidateUpdate(request)
}
