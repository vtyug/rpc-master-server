package validator

import (
	"log"
	"regexp"
	"unicode"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	mobileRegex       = regexp.MustCompile(`^1[3-9]\d{9}$`)
	emailRegex        = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	usernameRegex     = regexp.MustCompile(`^[a-zA-Z0-9_]{4,16}$`)
	passwordRegex     = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*]{6,20}$`)
	urlRegex          = regexp.MustCompile(`^(http|https):\/\/[^\s$.?#].[^\s]*$`)
	ipRegex           = regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
	dateRegex         = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	postalCodeRegex   = regexp.MustCompile(`^\d{6}$`)
	idCardRegex       = regexp.MustCompile(`^\d{15}|\d{18}$`)
	idRegex           = regexp.MustCompile(`^[1-9]\d*$`)
	apiKeyRegex       = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
	nameRegex         = regexp.MustCompile(`^\S+$`)
	collectionIDRegex = regexp.MustCompile(`^[1-9]\d*$`)
	parentIDRegex     = regexp.MustCompile(`^[1-9]\d*$`)
)

func Setup() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		log.Println("Registering custom validators")
		registerCustomValidators(v)
	} else {
		log.Println("Failed to register custom validators")
	}
}

func registerCustomValidators(v *validator.Validate) {
	validations := map[string]func(validator.FieldLevel) bool{
		"mobile":        validateMobile,
		"email":         validateEmail,
		"username":      validateUsername,
		"password":      validatePassword,
		"chinese":       validateChinese,
		"url":           validateURL,
		"ip":            validateIP,
		"date":          validateDate,
		"postalcode":    validatePostalCode,
		"idcard":        validateIDCard,
		"id":            validateID,
		"api_key":       validateAPIKey,
		"name":          validateName,
		"collection_id": validateCollectionID,
		"parent_id":     validateParentID,
	}

	for tag, fn := range validations {
		if err := v.RegisterValidation(tag, fn); err != nil {
			log.Printf("Failed to register validation for %s: %v", tag, err)
		}
	}
}

func validateMobile(fl validator.FieldLevel) bool {
	return mobileRegex.MatchString(fl.Field().String())
}

func validateEmail(fl validator.FieldLevel) bool {
	return emailRegex.MatchString(fl.Field().String())
}

func validateUsername(fl validator.FieldLevel) bool {
	return usernameRegex.MatchString(fl.Field().String())
}

func validatePassword(fl validator.FieldLevel) bool {
	return passwordRegex.MatchString(fl.Field().String())
}

func validateChinese(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	for _, r := range str {
		if !unicode.Is(unicode.Han, r) {
			return false
		}
	}
	return true
}

func validateURL(fl validator.FieldLevel) bool {
	return urlRegex.MatchString(fl.Field().String())
}

func validateIP(fl validator.FieldLevel) bool {
	return ipRegex.MatchString(fl.Field().String())
}

func validateDate(fl validator.FieldLevel) bool {
	return dateRegex.MatchString(fl.Field().String())
}

func validatePostalCode(fl validator.FieldLevel) bool {
	return postalCodeRegex.MatchString(fl.Field().String())
}

func validateIDCard(fl validator.FieldLevel) bool {
	return idCardRegex.MatchString(fl.Field().String())
}

func validateID(fl validator.FieldLevel) bool {
	return idRegex.MatchString(fl.Field().String())
}

func validateAPIKey(fl validator.FieldLevel) bool {
	return apiKeyRegex.MatchString(fl.Field().String())
}

func validateName(fl validator.FieldLevel) bool {
	return nameRegex.MatchString(fl.Field().String())
}

func validateCollectionID(fl validator.FieldLevel) bool {
	id := fl.Field().String()
	return collectionIDRegex.MatchString(id)
}

func validateParentID(fl validator.FieldLevel) bool {
	id := fl.Field().String()
	return parentIDRegex.MatchString(id)
}

// TranslateError 翻译验证错为中文
func TranslateError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		messages := make([]string, 0, len(errs))
		for _, e := range errs {
			var msg string
			switch e.Tag() {
			case "required":
				msg = "字段 " + e.Field() + " 是必填项"
			case "username":
				msg = "用户名必须是4到16位的字母、数字或下划线"
			case "password":
				msg = "密码必须是6到20位的字母、数字或特殊字符"
			case "email":
				msg = "邮箱格式不正确"
			case "mobile":
				msg = "手机号格式不正确"
			case "url":
				msg = "URL 格式不正确"
			case "ip":
				msg = "IP 地址格式不正确"
			case "date":
				msg = "日期格式不正确，格式应为 YYYY-MM-DD"
			case "postalcode":
				msg = "邮政编码格式不正确"
			case "idcard":
				msg = "身份证号码格式不正确"
			case "id":
				msg = "ID 格式不正确，必须是正整数"
			case "api_key":
				msg = "API Key 格式不正确"
			case "name":
				msg = "名称不能为空且不能包含空格"
			case "collection_id":
				msg = "集合ID格式不正确，必须是正整数"
			case "parent_id":
				msg = "父文件夹ID格式不正确，必须是正整数"
			default:
				msg = "字段 " + e.Field() + " 验证失败"
			}
			messages = append(messages, msg)
		}
		return joinMessages(messages)
	}
	return err.Error()
}

func joinMessages(messages []string) string {
	if len(messages) == 0 {
		return ""
	}
	return messages[0] // 只返回第一个错误信息
}
