package http

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/my-org/base-go/internal/order/app"
)

type placeOrderRequest struct {
	CustomerID uuid.UUID           `json:"customer_id"`
	Currency   string              `json:"currency"`
	Items      []placeOrderItemReq `json:"items"`
}

type placeOrderItemReq struct {
	ProductID uuid.UUID `json:"product_id"`
	UnitPrice int64     `json:"unit_price"`
	Quantity  int       `json:"quantity"`
}

// Validate kiểm tra định dạng input. Quy tắc nghiệp vụ (đơn rỗng, cùng loại tiền...) do domain kiểm tra.
func (r placeOrderRequest) Validate() error {
	if r.CustomerID == uuid.Nil {
		return errors.New("customer_id is required")
	}
	if r.Currency != "" && len(r.Currency) != 3 {
		return errors.New("currency must be a 3-letter ISO code")
	}
	for i, it := range r.Items {
		if it.ProductID == uuid.Nil {
			return fmt.Errorf("items[%d].product_id is required", i)
		}
	}
	return nil
}

func (r placeOrderRequest) toCmd() app.PlaceOrderCmd {
	cmd := app.PlaceOrderCmd{CustomerID: r.CustomerID, Currency: r.Currency}
	for _, it := range r.Items {
		cmd.Items = append(cmd.Items, app.PlaceOrderItem{ProductID: it.ProductID, UnitPrice: it.UnitPrice, Quantity: it.Quantity})
	}
	return cmd
}
