package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type OrderID uuid.UUID

func (id OrderID) String() string { return uuid.UUID(id).String() }

type OrderItem struct {
	ProductID uuid.UUID
	UnitPrice Money
	Quantity  int
}

func (i OrderItem) Subtotal() Money { return i.UnitPrice.Mul(i.Quantity) }

type Order struct {
	id         OrderID
	customerID uuid.UUID
	items      []OrderItem
	status     Status
	events     []any
}

func NewOrder(customerID uuid.UUID, items []OrderItem) (*Order, error) {
	if len(items) == 0 {
		return nil, ErrEmptyOrder
	}
	currency := items[0].UnitPrice.Currency
	for _, it := range items {
		if it.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}
		if it.UnitPrice.Amount < 0 {
			return nil, ErrNegativePrice
		}
		if it.UnitPrice.Currency != currency {
			return nil, ErrMixedCurrency
		}
	}
	o := &Order{
		id:         OrderID(uuid.New()),
		customerID: customerID,
		items:      items,
		status:     StatusPending,
	}
	o.events = append(o.events, OrderPlaced{OrderID: o.id})
	return o, nil
}

// Rehydrate dựng lại Order từ DB, không phát event.
func Rehydrate(id OrderID, customerID uuid.UUID, items []OrderItem, st Status) *Order {
	return &Order{id: id, customerID: customerID, items: items, status: st}
}

func (o *Order) Cancel() error {
	if o.status == StatusShipped {
		return fmt.Errorf("%w: order %s already shipped", ErrNotCancellable, o.id)
	}
	o.status = StatusCancelled
	return nil
}

func (o *Order) Total() Money {
	total := Money{Currency: o.items[0].UnitPrice.Currency}
	for _, it := range o.items {
		total = total.Add(it.Subtotal())
	}
	return total
}

func (o *Order) ID() OrderID           { return o.id }
func (o *Order) CustomerID() uuid.UUID { return o.customerID }
func (o *Order) Items() []OrderItem    { return o.items }
func (o *Order) Status() Status        { return o.status }

func (o *Order) PullEvents() []any {
	ev := o.events
	o.events = nil
	return ev
}
