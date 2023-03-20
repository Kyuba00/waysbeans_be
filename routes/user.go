package routes

import (
	"waysbeans_be/handlers"
	"waysbeans_be/pkg/middleware"
	"waysbeans_be/pkg/mysql"
	"waysbeans_be/repositories"

	"github.com/labstack/echo/v4"
)

func UserRoutes(r *echo.Group) {
	userRepository := repositories.RepositoryUser(mysql.DB)
	h := handlers.HandlerUser(userRepository)
	r.GET("/users", h.FindUsers, middleware.Auth)
	r.GET("/user/:id", h.GetUser, middleware.Auth)
	r.POST("/user", h.CreateUser, middleware.Auth)
	r.PATCH("/user/:id", h.UpdateUser, middleware.Auth)
	r.DELETE("/user/:id", h.DeleteUser, middleware.Auth)
}
