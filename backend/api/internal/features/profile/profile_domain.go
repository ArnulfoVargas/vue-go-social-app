package profile

type ProfileService interface {
	GetProfile(userID string) (*Profile, error)
}
