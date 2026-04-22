package dto

type RegisterRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=50"`
	FullName    string `json:"full_name" validate:"required,min=1,max=100"`
	Email       string `json:"email" validate:"required,email,max=100"`
	PhoneNumber string `json:"phone_number" validate:"required,min=8,max=20"`
	Password    string `json:"password" validate:"required,min=6,max=100"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
	DeviceUUID string `json:"device_uuid,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
	DeviceType string `json:"device_type,omitempty"`
	Platform   string `json:"platform,omitempty"`
	AppVersion string `json:"app_version,omitempty"`
	OSVersion  string `json:"os_version,omitempty"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	User         UserResponse `json:"user"`
}
type ChangePasswordRequest struct {
	OldPassword        string `json:"old_password" validate:"required,min=6,max=100"`
	NewPassword        string `json:"new_password" validate:"required,min=6,max=100"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"required,min=6,max=100"`
}