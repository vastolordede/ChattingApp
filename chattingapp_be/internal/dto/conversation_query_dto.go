package dto

type ListConversationsRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}