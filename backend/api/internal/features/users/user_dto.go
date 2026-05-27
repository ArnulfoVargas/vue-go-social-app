package users

type UpdateProfileRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=3,max=50"`
	LastName *string `json:"lastname,omitempty" validate:"omitempty,min=3,max=50"`
	ImageUrl *string `json:"image,omitempty" validate:"omitempty,url"`
	Bio      *string `json:"bio,omitempty" validate:"omitempty,min=3,max=500"`
}

type UpdatedProfileResponse struct {
	User    User   `json:"user"`
	Message string `json:"message"`
}
