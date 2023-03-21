package handlers

import (
	"context"
	"net/http"
	"os"
	"time"
	dto "waysbeans_be/dto/result"
	userdto "waysbeans_be/dto/users"
	"waysbeans_be/pkg/bcrypt"
	"waysbeans_be/repositories"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type handlerUser struct {
	UserRepository repositories.UserRepository
}

func HandlerUser(UserRepository repositories.UserRepository) *handlerUser {
	return &handlerUser{UserRepository}
}

func (h *handlerUser) GetUser(c echo.Context) error {
	userInfo := c.Get("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	user, err := h.UserRepository.GetUser(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	var Image string
	if user.Image != "" {
		Image = os.Getenv("PATH_FILE") + user.Image
	}

	userResponse := userdto.UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Image:   Image,
		Phone:   user.Phone,
		Address: user.Address,
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: "Success", Data: userResponse})
}

// func (h *handlerUser) CreateUser(c echo.Context) error {
// 	request := new(userdto.CreateUserRequest)
// 	if err := c.Bind(request); err != nil {
// 		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
// 	}

// 	validation := validator.New()
// 	err := validation.Struct(request)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
// 	}

// 	user := models.User{
// 		Name:     request.Name,
// 		Email:    request.Email,
// 		Password: request.Password,
// 	}

// 	data, err := h.UserRepository.CreateUser(user)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
// 	}

// 	return c.JSON(http.StatusOK, dto.SuccessResult{Code: "Success", Data: data})
// }

func (h *handlerUser) UpdateUser(c echo.Context) error {
	// request := new(userdto.UpdateUserRequest)
	// if err := c.Bind(request); err != nil {
	// 	return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	// }

	// id, _ := strconv.Atoi(c.Param("id"))

	userInfo := c.Get("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	filePath := c.Get("dataFile").(string)

	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// UPLOAD FILE TO CLOUDINARY
	respImage, _ := cld.Upload.Upload(ctx, filePath, uploader.UploadParams{Folder: "uploads"})

	var requestImage string
	if respImage == nil {
		requestImage = filePath
	} else {
		requestImage = respImage.SecureURL
	}

	updateUser := userdto.UpdateUserRequest{
		Name:     c.FormValue("name"),
		Password: c.FormValue("password"),
		Image:    requestImage,
		Phone:    c.FormValue("phone"),
		Address:  c.FormValue("address"),
	}

	user, err := h.UserRepository.GetUser(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if updateUser.Name != "" {
		user.Name = updateUser.Name
	}

	// if updateUser.Email != "" {
	// 	user.Email = updateUser.Email
	// }

	if updateUser.Password != "" {
		password, _ := bcrypt.HashingPassword(updateUser.Password)
		user.Password = password
	}

	if filePath != "false" {
		user.Image = updateUser.Image
	}

	if updateUser.Address != "" {
		user.Address = updateUser.Address
	}

	if updateUser.Phone != "" {
		user.Phone = updateUser.Phone
	}

	user.UpdateAt = time.Now()

	data, err := h.UserRepository.UpdateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	updateResponse := userdto.UpdateUserResponse{
		Name:    data.Name,
		Image:   data.Image,
		Phone:   data.Phone,
		Address: data.Address,
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: "Success", Data: updateResponse})
}
