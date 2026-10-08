package app

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/my-org/base-go/internal/order/domain"
)

func existing(st domain.Status) *domain.Order {
	return domain.Rehydrate(domain.OrderID(uuid.New()), uuid.New(),
		[]domain.OrderItem{{ProductID: uuid.New(), UnitPrice: domain.Money{Amount: 10, Currency: "VND"}, Quantity: 1}}, st)
}

func TestCancelOrder(t *testing.T) {
	o := existing(domain.StatusPending)
	repo := newFakeRepo(o)
	if err := NewCancelOrder(repo).Execute(context.Background(), o.ID()); err != nil {
		t.Fatal(err)
	}
	if repo.orders[o.ID()].Status() != domain.StatusCancelled || repo.saves != 1 {
		t.Fatal("order not cancelled")
	}
}

func TestCancelOrder_NotFound(t *testing.T) {
	err := NewCancelOrder(newFakeRepo()).Execute(context.Background(), domain.OrderID(uuid.New()))
	if !errors.Is(err, ErrOrderNotFound) {
		t.Fatalf("err = %v", err)
	}
}

func TestCancelOrder_Shipped(t *testing.T) {
	o := existing(domain.StatusShipped)
	repo := newFakeRepo(o)
	if err := NewCancelOrder(repo).Execute(context.Background(), o.ID()); !errors.Is(err, domain.ErrNotCancellable) {
		t.Fatalf("err = %v", err)
	}
	if repo.saves != 0 {
		t.Fatal("must not save")
	}
}

func TestGetOrder(t *testing.T) {
	o := existing(domain.StatusPaid)
	v, err := NewGetOrder(newFakeRepo(o)).Execute(context.Background(), o.ID())
	if err != nil || v.Status != "paid" || v.Total != 10 || len(v.Items) != 1 {
		t.Fatalf("view=%+v err=%v", v, err)
	}
	if _, err := NewGetOrder(newFakeRepo()).Execute(context.Background(), o.ID()); !errors.Is(err, ErrOrderNotFound) {
		t.Fatalf("err = %v", err)
	}
}
