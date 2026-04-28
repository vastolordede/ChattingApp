package dto

type UserResponse struct {
	ID          int64   `json:"id"`
	Username    string  `json:"username"`
	FullName    string  `json:"full_name"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phone_number"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	Status      string  `json:"status"`
	IsVerified  bool    `json:"is_verified"`
	LastSeenAt  *string `json:"last_seen_at,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type UpdateProfileRequest struct {
	FullName  string  `json:"full_name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Bio       *string `json:"bio,omitempty"`
}

type SearchUsersRequest struct {
	Query string `json:"query"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

type UserSearchItem struct {
	ID         int64   `json:"id"`
	Username   string  `json:"username"`
	FullName   string  `json:"full_name"`
	AvatarURL  *string `json:"avatar_url,omitempty"`
	IsVerified bool    `json:"is_verified"`
	LastSeenAt *string `json:"last_seen_at,omitempty"`
}
