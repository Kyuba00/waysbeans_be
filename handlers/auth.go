package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"
	authdto "waysbeans_be/dto/auth"
	dto "waysbeans_be/dto/result"
	"waysbeans_be/models"
	"waysbeans_be/pkg/bcrypt"
	jwtToken "waysbeans_be/pkg/jwt"
	"waysbeans_be/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type handlerAuth struct {
	AuthRepository repositories.AuthRepository
}

func HandlerAuth(AuthRepository repositories.AuthRepository) *handlerAuth {
	return &handlerAuth{AuthRepository}
}

func (h *handlerAuth) Register(c echo.Context) error {
	request := new(authdto.RegisterRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	password, err := bcrypt.HashingPassword(request.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: password,
		Role:     "buyer",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	data, err := h.AuthRepository.Register(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	registerResponse := authdto.RegisterResponse{
		ID:    data.ID,
		Name:  data.Name,
		Email: data.Email,
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: "Success", Data: registerResponse})
}

func (h *handlerAuth) Login(c echo.Context) error {
	request := new(authdto.LoginRequest)
	if err := c.Bind(request); err != nil {
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Check email
	user, err := h.AuthRepository.Login(request.Email)
	if err != nil {
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		return c.JSON(http.StatusBadRequest, response)
	}

	// Check password
	isValid := bcrypt.CheckPasswordHash(request.Password, user.Password)
	if !isValid {
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "wrong email or password"}
		return c.JSON(http.StatusBadRequest, response)
	}

	//generate token
	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["role"] = user.Role                           //bisa digunakan untuk payload
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix() // 2 hours expired

	token, errGenerateToken := jwtToken.GenerateToken(&claims)
	if errGenerateToken != nil {
		log.Println(errGenerateToken)
		fmt.Println("Unauthorize")
		return c.JSON(http.StatusInternalServerError, errGenerateToken.Error())
	}

	loginResponse := authdto.LoginResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
		Token: token,
	}

	response := dto.SuccessResult{Code: "Success", Data: loginResponse}
	return c.JSON(http.StatusOK, response)
}

func (h *handlerAuth) CheckAuth(c echo.Context) error {
	userInfo := c.Get("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	//Check User by Id
	user, err := h.AuthRepository.GetUserAuth(userId)
	if err != nil {
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		return c.JSON(http.StatusBadRequest, response)
	}

	CheckAuthResponse := authdto.CheckAuthResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	response := dto.SuccessResult{Code: "Success", Data: CheckAuthResponse}
	return c.JSON(http.StatusOK, response)
}
