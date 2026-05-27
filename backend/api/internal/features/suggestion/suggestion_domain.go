package suggestion

import "Server/internal/features/users"

type SuggestionService interface {
	GetSuggestedUsers(id string) ([]users.User, error)
}
