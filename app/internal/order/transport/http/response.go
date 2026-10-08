package http

import (
	"github.com/google/uuid"

	"github.com/my-org/base-go/internal/order/app"
)

type orderCreatedResponse struct {
	ID uuid.UUID `json:"id"`
}

type orderResponse struct {
	ID         uuid.UUID           `json:"id"`
	CustomerID uuid.UUID           `json:"customer_id"`
	Status     string              `json:"status"`
	Currency   string              `json:"currency"`
	Total      int64               `json:"total"`
	Items      []orderItemResponse `json:"items"`
}

type orderItemResponse struct {
	ProductID uuid.UUID `json:"product_id"`
	UnitPrice int64     `json:"unit_price"`
	Quantity  int       `json:"quantity"`
}

func fromView(v app.OrderView) orderResponse {
	res := orderResponse{
		ID: v.ID, CustomerID: v.CustomerID, Status: v.Status, Currency: v.Currency, Total: v.Total,
		Items: make([]orderItemResponse, 0, len(v.Items)),
	}
	for _, it := range v.Items {
		res.Items = append(res.Items, orderItemResponse(it))
	}
	return res
}
