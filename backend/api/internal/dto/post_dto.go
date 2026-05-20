package dto

type PostRequest struct {
	Content string `json:"content,omitempty" validate:"max=500,omitempty"`
}
