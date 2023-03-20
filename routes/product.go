package routes

import (
	"waysbeans_be/handlers"
	"waysbeans_be/pkg/middleware"
	"waysbeans_be/pkg/mysql"
	"waysbeans_be/repositories"

	"github.com/labstack/echo/v4"
)

func ProductRoutes(e *echo.Group) {
	productRepository := repositories.RepositoryProduct(mysql.DB)
	h := handlers.HandlerProduct(productRepository)

	e.GET("/products", h.FindProducts, middleware.Auth)
	e.GET("/product/:id", h.GetProduct, middleware.Auth)
	e.POST("/product", h.CreateProduct, middleware.Auth, middleware.UploadFile)      //memerlukan UploadFile untuk bisa mengupdate gambar
	e.PATCH("/product/:id", h.UpdateProduct, middleware.Auth, middleware.UploadFile) //memerlukan UploadFile untuk bisa memasukkan gambar
	e.DELETE("/product/:id", h.DeleteProduct, middleware.Auth)

}
