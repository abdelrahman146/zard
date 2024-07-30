package userapi

type LoginWithEmailAndPasswordRequest struct {
	Email    string `json:"email,omitempty" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
}

type GenerateOTPForEmailRequest struct {
	Email string `json:"email,omitempty" validate:"required,email"`
}

type VerifyOTPForEmailRequest struct {
	Value string `json:"value,omitempty" validate:"required"`
	Otp   string `json:"otp,omitempty" validate:"required"`
}

// Implement User Context
