package handlers

import (
	"net/http"
	"strconv"
	"time"
	cartdto "waysbeans_be/dto/cart"
	dto "waysbeans_be/dto/result"
	"waysbeans_be/models"
	"waysbeans_be/repositories"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type handlerCart struct {
	CartRepository repositories.CartRepository
}

func HandlerCart(CartRepository repositories.CartRepository) *handlerCart {
	return &handlerCart{CartRepository}
}

func (h *handlerCart) CreateCart(c echo.Context) error {
	// GET USER ROLE FROM JWT TOKEN
	userInfo := c.Get("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	// GET REQUEST AND DECODING JSON
	cartRequest := new(cartdto.CartRequest)
	if err := c.Bind(cartRequest); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	// RUN REPOSITORY GET PRODUCT BY PRODUCT ID
	product, err := h.CartRepository.GetProductCartByID(int(cartRequest.ProductID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	// FIND TOTAL PRICE PRODUCT FROM QUANTITY REQUEST
	total := product.Price * cartRequest.Quantity

	// RUN REPOSITORY GET TRANSACTION BY USER ID
	userTransaction, err := h.CartRepository.GetCartTransactionByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	// CHECK IF EXIST
	if userTransaction.ID == 0 {
		// SETUP FOR QUERY TRANSACTION
		transaction := models.Transaction{
			ID:       int(time.Now().Unix()),
			UserID:   userID,
			Status:   "waiting",
			Total:    0,
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
		}

		// RUN REPOSITORY CREATE TRANSACTION
		transactionData, err := h.CartRepository.CreateTransaction(transaction)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
		}

		// SETUP FOR QUERY CART
		cart := models.Cart{
			UserID:        userID,
			ProductID:     cartRequest.ProductID,
			Product:       models.Product{},
			OrderQty:      cartRequest.Quantity,
			Subtotal:      total,
			TransactionID: transactionData.ID,
			CreateAt:      time.Now(),
		}

		// RUN REPOSITORY CREATE CART
		data, err := h.CartRepository.CreateCart(cart)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
		}

		dataResponse, _ := h.CartRepository.GetCart(int(data.ID))

		// WRITE RESPONSE
		return c.JSON(http.StatusOK, dto.SuccessResult{Code: "success", Data: dataResponse})
	} else {
		// SETUP FOR QUERY CART
		cart := models.Cart{
			UserID:        userID,
			ProductID:     cartRequest.ProductID,
			Product:       models.Product{},
			OrderQty:      cartRequest.Quantity,
			Subtotal:      total,
			TransactionID: userTransaction.ID,
			CreateAt:      time.Now(),
		}

		// RUN REPOSITORY CREATE CART
		data, err := h.CartRepository.CreateCart(cart)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
		}

		dataResponse, _ := h.CartRepository.GetCart(int(data.ID))

		// WRITE RESPONSE
		return c.JSON(http.StatusOK, dto.SuccessResult{Code: "success", Data: dataResponse})
	}
}

func (h *handlerCart) DeleteCart(c echo.Context) error {
	// GET CART ID FROM URL
	cartID, _ := strconv.Atoi(c.Param("id"))

	// GET USER ROLE FROM JWT TOKEN
	userInfo := c.Get("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	// GET CART
	cart, err := h.CartRepository.GetCart(cartID)
	if err != nil {
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		return c.JSON(http.StatusBadRequest, response)
	}

	// VALIDATE REQUEST BY USER ID
	if userID != int(cart.UserID) {
		response := dto.ErrorResult{Code: http.StatusUnauthorized, Message: "unauthorized"}
		return c.JSON(http.StatusUnauthorized, response)
	}

	// DELETE DATA
	data, err := h.CartRepository.DeleteCart(cart)
	if err != nil {
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		return c.JSON(http.StatusInternalServerError, response)
	}

	// WRITE RESPONSE
	response := dto.SuccessResult{Code: "success", Data: data}
	return c.JSON(http.StatusOK, response)
}

func (h *handlerCart) FindCartByTransactionID(c echo.Context) error {
	// GET USER ROLE FROM JWT TOKEN
	userInfo := c.Get("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	// RUN REPOSITORY GET TRANSACTION BY USER ID
	transaction, err := h.CartRepository.GetCartTransactionByUserID(userID)
	if err != nil {
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		return c.JSON(http.StatusBadRequest, response)
	}

	// RUN REPOSITORY FIND CARTS BY TRANSACTION ID
	carts, err := h.CartRepository.FindCartByTransactionID(int(transaction.ID))
	if err != nil {
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		return c.JSON(http.StatusBadRequest, response)
	}

	// WRITE RESPONSE
	response := dto.SuccessResult{Code: "success", Data: carts}
	return c.JSON(http.StatusOK, response)
}
