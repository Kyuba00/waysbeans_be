package cartdto

type CartRequest struct {
	ProductID int `json:"id"`
	Quantity  int `json:"orderQuantity"`
}
