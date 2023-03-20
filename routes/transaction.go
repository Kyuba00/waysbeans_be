package routes

import (
	"waysbeans_be/handlers"
	"waysbeans_be/pkg/middleware"
	"waysbeans_be/pkg/mysql"
	"waysbeans_be/repositories"

	"github.com/labstack/echo/v4"
)

func TransactionRoutes(r *echo.Group) {
	transactionRepository := repositories.RepositoryTransaction(mysql.DB)
	h := handlers.HandlerTransaction(transactionRepository)

	r.GET("/transactions", h.FindTransactions, middleware.Auth)
	r.GET("/transaction-id", h.GetTransaction, middleware.Auth)
	r.POST("/transaction", h.CreateTransaction, middleware.Auth)
	r.DELETE("/transaction/:id", h.DeleteTransaction, middleware.Auth)
	r.PATCH("/transactionID", h.UpdateTransaction, middleware.Auth)
	r.POST("/notification", h.Notification)
	r.GET("/transaction-status", h.FindbyIDTransaction, middleware.Auth)
	r.GET("/transaction1", h.AllProductById, middleware.Auth)
}
