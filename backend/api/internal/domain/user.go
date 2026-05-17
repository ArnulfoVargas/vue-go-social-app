package domain

type UserService interface {
	GetUser(id string) (*User, error)
}

type UserRepository interface {
	GetUser(id string) (*User, error)
}
