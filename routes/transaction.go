package routes

import (
	"waysbeans_be/handlers"
	"waysbeans_be/pkg/middleware"
	"waysbeans_be/pkg/mysql"
	"waysbeans_be/repositories"

	"github.com/labstack/echo/v4"
)

func TransactionRoutes(e *echo.Group) {
	transactionRepository := repositories.RepositoryTransaction(mysql.DB)
	h := handlers.HandlerTransaction(transactionRepository)

	e.GET("/admin/transaction", h.FindTransactions, middleware.Auth)
	e.PATCH("/transaction", h.UpdateTransaction, middleware.Auth)
	e.GET("/user/transaction", h.GetUserTransactionByUserID, middleware.Auth)
	e.POST("/notification", h.Notification)
}
