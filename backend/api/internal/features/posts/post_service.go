package posts

import (
	"Server/internal/features/likes"
	"Server/internal/features/users"
	"Server/internal/helpers"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type postService struct {
	postRepo PostRepository
	likeRepo likes.LikeRepository
	userRepo users.UserRepository
}

func NewPostService(postRepo PostRepository, likeRepo likes.LikeRepository, userRepo users.UserRepository) *postService {
	return &postService{
		userRepo: userRepo,
		postRepo: postRepo,
		likeRepo: likeRepo,
	}
}

func (p *postService) CreatePost(userID string, post PostAdd) (Post, error) {
	userId, err := helpers.ToObjectID(userID)
	if err != nil {
		return Post{}, err
	}

	exists, err := p.userRepo.UserExistsById(userId)
	if err != nil {
		return Post{}, err
	}
	if !exists {
		return Post{}, errors.New("user not found")
	}

	postModel := Post{
		UserID:    userId,
		Content:   post.Content,
		Status:    1,
		Media:     post.Media,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	err = p.postRepo.CreatePost(postModel)
	if err != nil {
		return Post{}, err
	}

	return postModel, nil
}

func (p *postService) GetPost(postId string) (Post, error) {
	pId, err := helpers.ToObjectID(postId)
	if err != nil {
		return Post{}, err
	}
	return p.postRepo.GetPost(pId)
}

func (p *postService) DeletePost(postId string) error {
	pId, err := helpers.ToObjectID(postId)
	if err != nil {
		return err
	}

	err = p.postRepo.DeletePost(pId)
	if err != nil {
		return err
	}
	// TODO: delete likes associated with this post
	p.likeRepo.DeleteLikesFromPost(pId)

	// TODO: delete images associated with this post
	// TODO: delete comments associated with this post

	return nil
}

func (p *postService) UpdatePost(postId string, req UpdatePostRequest) (Post, error) {
	pId, err := helpers.ToObjectID(postId)
	if err != nil {
		return Post{}, err
	}

	if req.Content == "" {
		return Post{}, errors.New("content is required")
	}

	post := bson.M{
		"content":   req.Content,
		"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
	}

	p.postRepo.UpdatePost(pId, post)

	return Post{}, nil
}

func (p *postService) GetPostsByUserId(userId string) ([]Post, error) {
	uId, err := helpers.ToObjectID(userId)
	if err != nil {
		return nil, err
	}

	return p.postRepo.GetPostsByUserId(uId)
}

func (p *postService) ToggleLike(postId string, userId string) error {
	pId, err := helpers.ToObjectID(postId)
	if err != nil {
		return err
	}
	uId, err := helpers.ToObjectID(userId)
	if err != nil {
		return err
	}

	exists, err := p.userRepo.UserExistsById(uId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user does not exist")
	}

	exists, err = p.postRepo.ExistsById(pId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("post does not exist")
	}

	hasLike, err := p.likeRepo.HasLike(pId, uId)
	if err != nil {
		return err
	}
	if hasLike {
		err = p.likeRepo.DeleteLike(likes.Like{PostID: pId, UserID: uId})
		if err != nil {
			return err
		}
		return nil
	}
	like := likes.Like{
		PostID:    pId,
		UserID:    uId,
		ID:        primitive.NewObjectID(),
		Status:    1,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	err = p.likeRepo.AddLike(like)
	if err != nil {
		return err
	}

	return nil
}

func (p *postService) GetSuggestedPosts(userId string, limit int) ([]Post, error) {
	uId, err := helpers.ToObjectID(userId)
	if err != nil {
		return nil, err
	}

	return p.postRepo.GetSuggestedPosts(uId, limit)
}
