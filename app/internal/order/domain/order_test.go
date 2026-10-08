package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func vnd(a int64) Money { return Money{Amount: a, Currency: "VND"} }

func TestNewOrder(t *testing.T) {
	pid := uuid.New()
	tests := []struct {
		name    string
		items   []OrderItem
		wantErr error
	}{
		{"ok", []OrderItem{{pid, vnd(100), 2}}, nil},
		{"empty", nil, ErrEmptyOrder},
		{"zero qty", []OrderItem{{pid, vnd(100), 0}}, ErrInvalidQuantity},
		{"negative price", []OrderItem{{pid, vnd(-1), 1}}, ErrNegativePrice},
		{"mixed currency", []OrderItem{{pid, vnd(1), 1}, {pid, Money{1, "USD"}, 1}}, ErrMixedCurrency},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := NewOrder(uuid.New(), tt.items)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if err == nil && o.Status() != StatusPending {
				t.Fatalf("status = %s", o.Status())
			}
		})
	}
}

func TestTotal(t *testing.T) {
	o, _ := NewOrder(uuid.New(), []OrderItem{
		{uuid.New(), vnd(100), 2},
		{uuid.New(), vnd(50), 3},
	})
	if got := o.Total(); got != vnd(350) {
		t.Fatalf("total = %+v", got)
	}
}

func TestPullEvents(t *testing.T) {
	o, _ := NewOrder(uuid.New(), []OrderItem{{uuid.New(), vnd(1), 1}})
	ev := o.PullEvents()
	if len(ev) != 1 || ev[0].(OrderPlaced).OrderID != o.ID() {
		t.Fatalf("events = %v", ev)
	}
	if len(o.PullEvents()) != 0 {
		t.Fatal("events must be cleared after pull")
	}
}

func TestCancel(t *testing.T) {
	o := Rehydrate(OrderID(uuid.New()), uuid.New(), []OrderItem{{uuid.New(), vnd(1), 1}}, StatusPaid)
	if err := o.Cancel(); err != nil || o.Status() != StatusCancelled {
		t.Fatalf("cancel paid: err=%v status=%s", err, o.Status())
	}

	shipped := Rehydrate(OrderID(uuid.New()), uuid.New(), []OrderItem{{uuid.New(), vnd(1), 1}}, StatusShipped)
	if err := shipped.Cancel(); !errors.Is(err, ErrNotCancellable) {
		t.Fatalf("cancel shipped: err=%v", err)
	}
	if len(shipped.PullEvents()) != 0 {
		t.Fatal("rehydrated order must not emit events")
	}
}
