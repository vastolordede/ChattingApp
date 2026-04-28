package dto

type RegisterDeviceRequest struct {
	DeviceUUID string `json:"device_uuid" validate:"required"`
	DeviceName string `json:"device_name" validate:"required"`
	DeviceType string `json:"device_type" validate:"required"` // phone/tablet/desktop/web
	Platform   string `json:"platform" validate:"required"`    // android/ios/windows/macos/linux/web
	AppVersion string `json:"app_version,omitempty"`
	OSVersion  string `json:"os_version,omitempty"`
	PushToken  string `json:"push_token,omitempty"`
}

type UpdatePushTokenRequest struct {
	PushToken string `json:"push_token" validate:"required"`
}

type UserDeviceResponse struct {
	ID           int64   `json:"id"`
	UserID       int64   `json:"user_id"`
	DeviceUUID   string  `json:"device_uuid"`
	DeviceName   string  `json:"device_name"`
	DeviceType   string  `json:"device_type"`
	Platform     string  `json:"platform"`
	AppVersion   *string `json:"app_version,omitempty"`
	OSVersion    *string `json:"os_version,omitempty"`
	PushToken    *string `json:"push_token,omitempty"`
	IsTrusted    bool    `json:"is_trusted"`
	IsActive     bool    `json:"is_active"`
	LastSeenAt   *string `json:"last_seen_at,omitempty"`
	RegisteredAt string  `json:"registered_at"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}
