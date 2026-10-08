package app

import (
	"context"
	"fmt"

	"github.com/my-org/base-go/internal/order/domain"
)

type GetOrder struct {
	repo Repository
}

func NewGetOrder(r Repository) *GetOrder { return &GetOrder{repo: r} }

func (uc *GetOrder) Execute(ctx context.Context, id domain.OrderID) (OrderView, error) {
	order, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return OrderView{}, fmt.Errorf("find order: %w", err)
	}
	if order == nil {
		return OrderView{}, ErrOrderNotFound
	}
	return toView(order), nil
}
