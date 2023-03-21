package userdto

type UserResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Image    string `json:"image"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type UpdateUserResponse struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type DeleteUserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
