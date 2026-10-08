package app

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/my-org/base-go/internal/order/domain"
)

type fakeRepo struct {
	orders map[domain.OrderID]*domain.Order
	saves  int
}

func newFakeRepo(orders ...*domain.Order) *fakeRepo {
	r := &fakeRepo{orders: map[domain.OrderID]*domain.Order{}}
	for _, o := range orders {
		r.orders[o.ID()] = o
	}
	return r
}

func (f *fakeRepo) Save(_ context.Context, o *domain.Order) error {
	f.orders[o.ID()] = o
	f.saves++
	return nil
}

func (f *fakeRepo) FindByID(_ context.Context, id domain.OrderID) (*domain.Order, error) {
	return f.orders[id], nil
}

type fakePayment struct {
	err     error
	charged []domain.Money
}

func (f *fakePayment) Charge(_ context.Context, _ uuid.UUID, m domain.Money) error {
	if f.err != nil {
		return f.err
	}
	f.charged = append(f.charged, m)
	return nil
}

type fakePublisher struct{ events []any }

func (f *fakePublisher) Publish(_ context.Context, ev []any) error {
	f.events = append(f.events, ev...)
	return nil
}

var errBoom = errors.New("boom")
