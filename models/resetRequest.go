package models

type ForgotPassword struct {
	Email string `json:"email" binding:"required"`
}

type ResetRequest struct {
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}
