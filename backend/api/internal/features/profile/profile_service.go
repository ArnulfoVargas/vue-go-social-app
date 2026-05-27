package profile

import (
	"Server/internal/features/follows"
	"Server/internal/features/media"
	"Server/internal/features/posts"
	"Server/internal/features/users"
	"Server/internal/helpers"
	"errors"
	"sync"
)

type profileService struct {
	userService   users.UserRepository
	followService follows.FollowRepository
	postService   posts.PostRepository
}

func NewProfileService(userRepo users.UserRepository, followRepo follows.FollowRepository, postRepo posts.PostRepository) *profileService {
	return &profileService{
		userService:   userRepo,
		followService: followRepo,
		postService:   postRepo,
	}
}

func (s *profileService) GetProfile(userID string) (*Profile, error) {
	uid, err := helpers.ToObjectID(userID)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var user *users.User
	userPosts := make([]posts.Post, 0)
	var followers, following int64 = 0, 0

	var errs []error
	var (
		e1 error
		e2 error
		e3 error
		e4 error
	)

	wg.Go(func() {
		user, e1 = s.userService.GetUserById(uid)
		if err != nil {
			errs = append(errs, e1)
		}
	})

	wg.Go(func() {
		userPosts, e2 = s.postService.GetPostsByUserId(uid)
		if e2 != nil {
			errs = append(errs, e2)
		}
	})

	wg.Go(func() {
		followers, e3 = s.followService.GetFollowersCount(uid)
		if e3 != nil {
			errs = append(errs, e3)
		}
	})

	wg.Go(func() {
		following, e4 = s.followService.GetFollowingCount(uid)
		if e4 != nil {
			errs = append(errs, e4)
		}
	})

	wg.Wait()

	if len(errs) > 0 {
		return nil, errors.New("cannot get profile")
	}

	profile := &Profile{
		ID:        user.ID.Hex(),
		Name:      user.Name,
		Posts:     posts.PostsFromModels(userPosts),
		Followers: followers,
		Following: following,
	}

	if user.Avatar.URL != "" {
		profile.Avatar = &media.MediaResponse{
			ID:  user.Avatar.ID.Hex(),
			URL: user.Avatar.URL,
		}
	}

	return profile, nil
}
