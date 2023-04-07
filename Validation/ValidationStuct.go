package validation

type SendOtp struct{
	Phone string `json:"phone" validate:"required,e164"`
}