package user

// RegisterUserDTO example
// @Description Payload to register a user
type RegisterUserDTO struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email" binding:"required,email"`
	OTP string `json:"otp"`
}

type SendOTPRequestDTO struct {
	Email   string `json:"email" binding:"required"`
	Subject string `json:"subject"`
}

type VerifyOTPRequestDTO struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type LoginRequestDTO struct {
	Email string `json:"email" binding:"required"`
	OTP   string `json:"otp"`
}
