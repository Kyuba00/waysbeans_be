package userdto

type CreateUserRequest struct {
	Name     string `json:"name" form:"name" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
	Role     string `json:"role" form:"role" validate:"required"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" form:"name"`
	Password string `json:"password" form:"password"`
	Image    string `json:"image" form:"image"`
	Phone    string `json:"phone" form:"phone"`
	Address  string `json:"address" form:"address"`
}
