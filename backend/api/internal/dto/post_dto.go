package dto

type PostDTO struct {
	ID      string `json:"id"`
	Content string `json:"content,omitempty"`
}
