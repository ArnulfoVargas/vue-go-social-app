package profile

import (
	"Server/internal/features/media"
	"Server/internal/features/posts"
)

type Profile struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	Avatar    *media.MediaResponse `json:"avatar,omitempty"`
	Posts     []posts.PostResponse `json:"posts"`
	Followers int64                `json:"followers"`
	Following int64                `json:"following"`
}
