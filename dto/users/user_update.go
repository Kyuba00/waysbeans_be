package usersdto

type UpdateUserRequest struct {
	Name  string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Status   string `json:"status" form:"status"`
}
