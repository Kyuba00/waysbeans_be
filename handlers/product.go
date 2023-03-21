package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	productdto "waysbeans_be/dto/product"
	dto "waysbeans_be/dto/result"
	"waysbeans_be/models"
	"waysbeans_be/repositories"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// var path_file = "http://localhost:5000/uploads/"

type handlerProduct struct {
	ProductRepository repositories.ProductRepository
}

func HandlerProduct(ProductRepository repositories.ProductRepository) *handlerProduct {
	return &handlerProduct{ProductRepository}
}

func (h *handlerProduct) FindProducts(c echo.Context) error {
	products, err := h.ProductRepository.FindProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	var productsResponse []productdto.ProductResponse
	for _, p := range products {
		productResponse := productdto.ProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Price:       p.Price,
			Description: p.Description,
			Image:       os.Getenv("PATH_FILE") + p.Image,
			Stock:       p.Stock,
		}
		productsResponse = append(productsResponse, productResponse)
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: "Success", Data: productsResponse})
}

func (h *handlerProduct) GetProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	product, err := h.ProductRepository.GetProduct(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	productResponse := productdto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Image:       os.Getenv("PATH_FILE") + product.Image,
		Stock:       product.Stock,
	}

	response := dto.SuccessResult{Code: "Success", Data: productResponse}
	return c.JSON(http.StatusOK, response)
}

func (h *handlerProduct) CreateProduct(c echo.Context) error {
	userInfo := c.Get("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]

	if userRole != "admin" {
		response := dto.ErrorResult{Code: http.StatusUnauthorized, Message: "Unauthorized"}
		return c.JSON(http.StatusUnauthorized, response)
	}

	filePath := c.Get("dataFile").(string)

	fmt.Println(filePath)

	// Declare Context Background, Cloud Name, API Key, API Secret ...
	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// Upload file to Cloudinary ...
	resp, err := cld.Upload.Upload(ctx, filePath, uploader.UploadParams{Folder: "waysbeans"})
	if err != nil {
		fmt.Println(err.Error())
	}

	price, _ := strconv.Atoi(c.FormValue("price"))
	stock, _ := strconv.Atoi(c.FormValue("stock"))

	request := productdto.ProductRequest{
		Name:        c.FormValue("name"),
		Price:       price,
		Image:       resp.SecureURL,
		Description: c.FormValue("description"),
		Stock:       stock,
	}

	validation := validator.New()
	err = validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	product := models.Product{
		Name:        request.Name,
		Price:       request.Price,
		Stock:       request.Stock,
		Image:       request.Image,
		Description: request.Description,
		CreateAt:    time.Now(),
		UpdateAt:    time.Now(),
	}

	product, err = h.ProductRepository.CreateProduct(product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	productResponse := productdto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Image:       product.Image,
		Stock:       product.Stock,
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: "Success", Data: productResponse})
}

func (h *handlerProduct) UpdateProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	userInfo := c.Get("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]

	if userRole != "admin" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResult{Code: http.StatusUnauthorized, Message: "Unauthorized"})
	}

	filePath := c.Get("dataFile").(string)

	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// Upload file to Cloudinary ...
	resp, err := cld.Upload.Upload(ctx, filePath, uploader.UploadParams{Folder: "waysbeans"})
	if err != nil {
		fmt.Println(err.Error())
	}

	price, _ := strconv.Atoi(c.FormValue("price"))
	stock, _ := strconv.Atoi(c.FormValue("stock"))
	request := productdto.UpdateProductRequest{
		Name:        c.FormValue("name"),
		Price:       price,
		Image:       resp.SecureURL,
		Description: c.FormValue("description"),
		Stock:       stock,
	}

	product, err := h.ProductRepository.GetProduct(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if request.Name != "" {
		product.Name = request.Name
	}

	if request.Price != 0 {
		product.Price = request.Price
	}

	if filePath != "false" {
		product.Image = request.Image
	}

	if request.Stock != 0 {
		product.Stock = request.Stock
	}

	if request.Description != "" {
		product.Description = request.Description
	}

	product.UpdateAt = time.Now()

	data, err := h.ProductRepository.UpdateProduct(product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	updateResponse := productdto.ProductResponse{
		ID:          data.ID,
		Name:        data.Name,
		Price:       data.Price,
		Description: data.Description,
		Image:       data.Image,
		Stock:       data.Stock,
	}

	response := dto.SuccessResult{Code: "Success", Data: updateResponse}
	return c.JSON(http.StatusOK, response)
}

func (h *handlerProduct) DeleteProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	userInfo := c.Get("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]

	if userRole != "admin" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResult{Code: http.StatusUnauthorized, Message: "Unauthorized"})
	}

	delete, err := h.ProductRepository.GetProduct(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	data, err := h.ProductRepository.DeleteProduct(delete)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	deleteResponse := productdto.DeleteProductResponse{
		ID:   data.ID,
		Name: data.Name,
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: "Success", Data: deleteResponse})
}
