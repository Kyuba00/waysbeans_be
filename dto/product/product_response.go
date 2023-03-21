package productdto

type ProductResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Stock       int    `json:"stock"`
}

type DeleteProductResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
