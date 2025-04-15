package user

// RegisterUserDTO example
// @Description Payload to register a user
type RegisterUserDTO struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	OTP      string `json:"otp"`
}

type LoginUserDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}


type SendOTPRequestDTO struct {
	Email   string `json:"email" binding:"required"`
	Subject string `json:"subject"`
}

type VerifyOTPRequestDTO struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}