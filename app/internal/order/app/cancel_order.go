package app

import (
	"context"
	"fmt"

	"github.com/my-org/base-go/internal/order/domain"
)

type CancelOrder struct {
	repo Repository
}

func NewCancelOrder(r Repository) *CancelOrder { return &CancelOrder{repo: r} }

func (uc *CancelOrder) Execute(ctx context.Context, id domain.OrderID) error {
	order, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find order: %w", err)
	}
	if order == nil {
		return ErrOrderNotFound
	}
	if err := order.Cancel(); err != nil {
		return err
	}
	if err := uc.repo.Save(ctx, order); err != nil {
		return fmt.Errorf("save order: %w", err)
	}
	return nil
}
