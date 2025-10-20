// ValidatorWrapper.go
package controllers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/gbrayhan/microservices-go/src/domain/constants"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	"github.com/go-playground/validator/v10"
)

// CommonValidator 通用验证器
// 用于验证各种数据结构的更新和创建操作
type CommonValidator struct {
	validate       *validator.Validate
	validationMap  map[string]string
	errorsMessages []string
}

// NewCommonValidator 创建新的通用验证器实例
// validationMap: 字段名到验证规则的映射
// 规则格式遵循 go-playground/validator 标准
func NewCommonValidator(validationMap map[string]string) *CommonValidator {
	v := &CommonValidator{
		validate:       validator.New(),
		validationMap:  validationMap,
		errorsMessages: make([]string, 0),
	}

	v.registerCustomValidations()
	return v
}

func (v *CommonValidator) registerCustomValidations() {
	_ = v.validate.RegisterValidation("update_validation", func(fl validator.FieldLevel) bool {
		m, ok := fl.Field().Interface().(map[string]any)
		if !ok {
			return false
		}

		for k, rule := range v.validationMap {
			if val, exists := m[k]; exists {
				errValidate := v.validate.Var(val, rule)
				if errValidate != nil {
					validatorErr := errValidate.(validator.ValidationErrors)
					v.errorsMessages = append(
						v.errorsMessages,
						fmt.Sprintf("%s does not satisfy condition %v=%v", k, validatorErr[0].Tag(), validatorErr[0].Param()),
					)
				}
			}
		}
		return true
	})

	// verify status
	_ = v.validate.RegisterValidation("status_enum", func(fl validator.FieldLevel) bool {
		value := fl.Field().Interface()
		switch v := value.(type) {
		case float64:
			// 转换为字符串并检查是否匹配枚举值
			strValue := fmt.Sprintf("%.0f", v)
			return strValue == constants.StatusEnabled || strValue == constants.StatusDisabled
		case int:
			strValue := fmt.Sprintf("%d", v)
			return strValue == constants.StatusEnabled || strValue == constants.StatusDisabled
		case string:
			return v == constants.StatusEnabled || v == constants.StatusDisabled
		default:
			return false
		}
	})

	// verify phone
	_ = v.validate.RegisterValidation("custom_phone", func(fl validator.FieldLevel) bool {
		phone, ok := fl.Field().Interface().(string)
		if !ok || phone == "" {
			return true // 允许空值由 omitempty 处理
		}
		// 自定义手机号正则表达式（示例为中国手机号）
		match := regexp.MustCompile(`^\+?\d{10,15}$`).MatchString(phone)
		return match
	})

}

func (v *CommonValidator) ValidateUpdate(request map[string]any) error {
	// 重置错误消息
	v.errorsMessages = make([]string, 0)

	// 执行自定义验证
	err := v.validate.Var(request, "update_validation")
	if err != nil {
		// 提供更友好的错误信息
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, validationErr := range validationErrors {
				field := validationErr.Field()
				tag := validationErr.Tag()
				param := validationErr.Param()

				errorMsg := v.getFriendlyErrorMessage(field, tag, param)
				v.errorsMessages = append(v.errorsMessages, errorMsg)
			}
		} else {
			v.errorsMessages = append(v.errorsMessages, "验证过程出错: "+err.Error())
		}
	}

	if len(v.errorsMessages) > 0 {
		return domainErrors.NewAppError(
			errors.New(strings.Join(v.errorsMessages, "; ")),
			domainErrors.ValidationError,
		)
	}

	return nil
}
func (v *CommonValidator) getFriendlyErrorMessage(field, tag, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("字段 '%s' 是必需的", field)
	case "max":
		return fmt.Sprintf("字段 '%s' 长度不能超过 %s 个字符", field, param)
	case "min":
		return fmt.Sprintf("字段 '%s' 长度不能少于 %s 个字符", field, param)
	case "email":
		return fmt.Sprintf("字段 '%s' 必须是有效的邮箱地址", field)
	// 添加更多友好错误信息
	default:
		return fmt.Sprintf("字段 '%s' 不满足条件 %s=%s", field, tag, param)
	}
}
