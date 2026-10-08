package app

import (
	"context"
	"fmt"

	"github.com/my-org/base-go/internal/order/domain"
)

const defaultCurrency = "VND"

type PlaceOrder struct {
	repo      Repository
	payment   PaymentGateway
	publisher Publisher
}

func NewPlaceOrder(r Repository, p PaymentGateway, pub Publisher) *PlaceOrder {
	return &PlaceOrder{repo: r, payment: p, publisher: pub}
}

func (uc *PlaceOrder) Execute(ctx context.Context, cmd PlaceOrderCmd) (domain.OrderID, error) {
	currency := cmd.Currency
	if currency == "" {
		currency = defaultCurrency
	}
	items := make([]domain.OrderItem, 0, len(cmd.Items))
	for _, it := range cmd.Items {
		items = append(items, domain.OrderItem{
			ProductID: it.ProductID,
			UnitPrice: domain.Money{Amount: it.UnitPrice, Currency: currency},
			Quantity:  it.Quantity,
		})
	}

	order, err := domain.NewOrder(cmd.CustomerID, items)
	if err != nil {
		return domain.OrderID{}, err
	}
	if err := uc.payment.Charge(ctx, cmd.CustomerID, order.Total()); err != nil {
		return domain.OrderID{}, fmt.Errorf("charge: %w", err)
	}
	if err := uc.repo.Save(ctx, order); err != nil {
		return domain.OrderID{}, fmt.Errorf("save order: %w", err)
	}
	if err := uc.publisher.Publish(ctx, order.PullEvents()); err != nil {
		return domain.OrderID{}, fmt.Errorf("publish events: %w", err)
	}
	return order.ID(), nil
}
