package domain

type LikeRepository interface {
	ToggleLike(postId string, userId string) error
}
