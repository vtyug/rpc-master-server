package validator

import (
	"log"
	"regexp"
	"unicode"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	mobileRegex     = regexp.MustCompile(`^1[3-9]\d{9}$`)
	emailRegex      = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	usernameRegex   = regexp.MustCompile(`^[a-zA-Z0-9_]{4,16}$`)
	passwordRegex   = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*]{6,20}$`)
	urlRegex        = regexp.MustCompile(`^(http|https):\/\/[^\s$.?#].[^\s]*$`)
	ipRegex         = regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
	dateRegex       = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	postalCodeRegex = regexp.MustCompile(`^\d{6}$`)
	idCardRegex     = regexp.MustCompile(`^\d{15}|\d{18}$`)
	idRegex         = regexp.MustCompile(`^[1-9]\d*$`)
	apiKeyRegex     = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
)

func Setup() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		log.Println("Registering custom validators")
		_ = v.RegisterValidation("mobile", validateMobile)
		_ = v.RegisterValidation("email", validateEmail)
		_ = v.RegisterValidation("username", validateUsername)
		_ = v.RegisterValidation("password", validatePassword)
		_ = v.RegisterValidation("chinese", validateChinese)
		_ = v.RegisterValidation("url", validateURL)
		_ = v.RegisterValidation("ip", validateIP)
		_ = v.RegisterValidation("date", validateDate)
		_ = v.RegisterValidation("postalcode", validatePostalCode)
		_ = v.RegisterValidation("idcard", validateIDCard)
		_ = v.RegisterValidation("id", validateID)
		_ = v.RegisterValidation("api_key", validateAPIKey)
	} else {
		log.Println("Failed to register custom validators")
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

// TranslateError 翻译验证错误为中文
func TranslateError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			switch e.Tag() {
			case "required":
				return "字段 " + e.Field() + " 是必填项"
			case "username":
				return "用户名必须是4到16位的字母、数字或下划线"
			case "password":
				return "密码必须是6到20位的字母、数字或特殊字符"
			case "email":
				return "邮箱格式不正确"
			case "mobile":
				return "手机号格式不正确"
			case "url":
				return "URL 格式不正确"
			case "ip":
				return "IP 地址格式不正确"
			case "date":
				return "日期格式不正确，格式应为 YYYY-MM-DD"
			case "postalcode":
				return "邮政编码格式不正确"
			case "idcard":
				return "身份证号码格式不正确"
			case "id":
				return "ID 格式不正确，必须是正整数"
			case "api_key":
				return "API Key 格式不正确"
			default:
				return "字段 " + e.Field() + " 验证失败"
			}
		}
	}
	return err.Error()
}
