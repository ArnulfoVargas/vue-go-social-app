package posts

import (
	"Server/internal/features/media"
)

type PostResponse struct {
	ID      string                `json:"id"`
	Content string                `json:"content,omitempty"`
	Media   []media.MediaResponse `json:"media,omitempty"`
}

type UpdatePostRequest struct {
	Content string `json:"content,omitempty" validate:"max=500,omitempty"`
}

type PostsResponse struct {
	Posts []PostResponse `json:"posts"`
}

type PostAdd struct {
	Content string        `json:"content,omitempty" validate:"max=500,omitempty"`
	Media   []media.Media `json:"media,omitempty"`
}

func PostsFromModels(posts []Post) []PostResponse {
	result := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		result = append(result, PostResponse{
			ID:      post.ID.Hex(),
			Content: post.Content,
			Media:   media.MediasFromModels(post.Media),
		})
	}
	return result
}
