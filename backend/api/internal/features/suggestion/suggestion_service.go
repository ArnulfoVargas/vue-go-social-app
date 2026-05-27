package suggestion

import (
	"Server/internal/constants"
	"Server/internal/features/follows"
	"Server/internal/features/users"
	"Server/internal/helpers"
)

type suggestionService struct {
	userRepo   users.UserRepository
	followRepo follows.FollowRepository
}

func NewSuggestionService(userRepo users.UserRepository, followRepo follows.FollowRepository) *suggestionService {
	return &suggestionService{
		userRepo:   userRepo,
		followRepo: followRepo,
	}
}

func (s *suggestionService) GetSuggestedUsers(id string) ([]users.User, error) {
	uId, err := helpers.ToObjectID(id)
	if err != nil {
		return nil, err
	}

	followingIds, err := s.followRepo.GetFollowingIds(uId)
	if err != nil {
		return nil, err
	}

	userId, err := helpers.ToObjectID(id)
	if err != nil {
		return nil, err
	}

	suggestedIds, err := s.followRepo.GetRelatedFollowSuggestions(userId, followingIds, constants.MAX_SUGGESTED_IDS)
	if err != nil {
		return nil, err
	}

	sugIdsLen := len(suggestedIds)
	if sugIdsLen < constants.MAX_SUGGESTED_IDS {
		excludedIds := append(followingIds, userId)
		excludedIds = append(excludedIds, suggestedIds...)

		randomUsers, err := s.userRepo.GetIdsExcluding(excludedIds, constants.MAX_SUGGESTED_IDS-sugIdsLen)

		if err != nil {
			return nil, err
		}

		return s.userRepo.GetUsersByIds(append(suggestedIds, randomUsers...))
	}

	return s.userRepo.GetUsersByIds(suggestedIds)
}
