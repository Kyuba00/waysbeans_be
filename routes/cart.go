package routes

import (
	"waysbeans_be/handlers"
	"waysbeans_be/pkg/middleware"
	"waysbeans_be/pkg/mysql"
	"waysbeans_be/repositories"

	"github.com/labstack/echo/v4"
)

func CartRoutes(r *echo.Group) {
	cartRepository := repositories.RepositoryCart(mysql.DB)
	h := handlers.HandlerCart(cartRepository)

	r.GET("/carts", h.FindCart)
	r.GET("/cart/:id", h.GetCart)
	r.POST("/cart", middleware.Auth(h.CreateCart))
	r.DELETE("/cart/:id", h.DeleteCart)
	r.PATCH("/cart/:id", middleware.Auth(h.UpdateCart))
}
