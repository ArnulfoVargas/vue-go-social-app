package users

import "Server/internal/features/media"

type UpdateProfileRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=3,max=50"`
	LastName *string `json:"lastname,omitempty" validate:"omitempty,min=3,max=50"`
	Bio      *string `json:"bio,omitempty" validate:"omitempty,min=3,max=500"`
}

type UpdatedProfileResponse struct {
	User    User   `json:"user"`
	Message string `json:"message"`
}

type SetProfilePictureResponse struct {
	AvatarUrl media.MediaResponse `json:"avatarUrl"`
}
