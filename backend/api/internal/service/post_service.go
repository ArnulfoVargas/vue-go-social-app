package service

import (
	"Server/internal/domain"
	"Server/internal/dto"
	"Server/internal/model"
)

type postService struct {
	postRepo domain.PostRepository
	likeRepo domain.LikeRepository
}

func NewPostService(postRepo domain.PostRepository, likeRepo domain.LikeRepository) *postService {
	return &postService{
		postRepo: postRepo,
		likeRepo: likeRepo,
	}
}

func (p *postService) CreatePost(userId string, post dto.PostRequest) (model.Post, error) {
	return model.Post{}, nil
}

func (p *postService) GetPost(postId string) (model.Post, error) {
	return model.Post{}, nil
}

func (p *postService) DeletePost(postId string) error {
	return nil
}

func (p *postService) UpdatePost(postId string, post dto.PostRequest) (model.Post, error) {
	return model.Post{}, nil
}

func (p *postService) GetPostsByUserId(userId string) ([]model.Post, error) {
	return nil, nil
}

func (p *postService) AttachImage(postId string, image *model.Media) error {
	return nil
}

func (p *postService) AttachManyImages(postId string, images []*model.Media) error {
	return nil
}

func (p *postService) ToggleLike(postId string, userId string) error {
	return nil
}

func (p *postService) DetachImage(postId string, imageId string) error {
	return nil
}

func (p *postService) DetachManyImages(postId string, imageIds []string) error {
	return nil
}
