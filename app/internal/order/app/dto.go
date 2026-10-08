package app

import (
	"github.com/google/uuid"

	"github.com/my-org/base-go/internal/order/domain"
)

type PlaceOrderCmd struct {
	CustomerID uuid.UUID
	Currency   string
	Items      []PlaceOrderItem
}

type PlaceOrderItem struct {
	ProductID uuid.UUID
	UnitPrice int64
	Quantity  int
}

// OrderView là DTO đọc, tách khỏi domain.Order để transport không phụ thuộc entity.
type OrderView struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	Status     string
	Currency   string
	Total      int64
	Items      []OrderItemView
}

type OrderItemView struct {
	ProductID uuid.UUID
	UnitPrice int64
	Quantity  int
}

func toView(o *domain.Order) OrderView {
	total := o.Total()
	v := OrderView{
		ID:         uuid.UUID(o.ID()),
		CustomerID: o.CustomerID(),
		Status:     string(o.Status()),
		Currency:   total.Currency,
		Total:      total.Amount,
		Items:      make([]OrderItemView, 0, len(o.Items())),
	}
	for _, it := range o.Items() {
		v.Items = append(v.Items, OrderItemView{ProductID: it.ProductID, UnitPrice: it.UnitPrice.Amount, Quantity: it.Quantity})
	}
	return v
}
