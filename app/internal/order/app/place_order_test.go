package app

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/my-org/base-go/internal/order/domain"
)

func validCmd() PlaceOrderCmd {
	return PlaceOrderCmd{
		CustomerID: uuid.New(),
		Items: []PlaceOrderItem{
			{ProductID: uuid.New(), UnitPrice: 100, Quantity: 2},
			{ProductID: uuid.New(), UnitPrice: 50, Quantity: 1},
		},
	}
}

func TestPlaceOrder_Success(t *testing.T) {
	repo, pay, pub := newFakeRepo(), &fakePayment{}, &fakePublisher{}
	id, err := NewPlaceOrder(repo, pay, pub).Execute(context.Background(), validCmd())
	if err != nil {
		t.Fatal(err)
	}
	if repo.orders[id] == nil {
		t.Fatal("order not saved")
	}
	if len(pay.charged) != 1 || pay.charged[0] != (domain.Money{Amount: 250, Currency: "VND"}) {
		t.Fatalf("charged = %+v", pay.charged)
	}
	if len(pub.events) != 1 {
		t.Fatalf("events = %v", pub.events)
	}
}

func TestPlaceOrder_DomainError(t *testing.T) {
	repo, pay := newFakeRepo(), &fakePayment{}
	_, err := NewPlaceOrder(repo, pay, &fakePublisher{}).Execute(context.Background(), PlaceOrderCmd{CustomerID: uuid.New()})
	if !errors.Is(err, domain.ErrEmptyOrder) {
		t.Fatalf("err = %v", err)
	}
	if len(pay.charged) != 0 || repo.saves != 0 {
		t.Fatal("must not charge or save an invalid order")
	}
}

func TestPlaceOrder_PaymentFails(t *testing.T) {
	repo := newFakeRepo()
	_, err := NewPlaceOrder(repo, &fakePayment{err: errBoom}, &fakePublisher{}).Execute(context.Background(), validCmd())
	if !errors.Is(err, errBoom) {
		t.Fatalf("err = %v", err)
	}
	if repo.saves != 0 {
		t.Fatal("must not save when payment fails")
	}
}
