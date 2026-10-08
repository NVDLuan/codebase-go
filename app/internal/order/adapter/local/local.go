// Package local chứa adapter không cần hạ tầng ngoài: dùng khi chạy local,
// hoặc khi chưa cấu hình DATABASE_URL / STRIPE_API_KEY / KAFKA_BROKERS.
// KHÔNG dùng cho production: dữ liệu mất khi pod restart, không thu tiền thật.
package local

import (
	"context"
	"log/slog"
	"sync"

	"github.com/google/uuid"

	"github.com/my-org/base-go/internal/order/domain"
)

// Repository lưu Order trong bộ nhớ.
type Repository struct {
	mu     sync.RWMutex
	orders map[domain.OrderID]*domain.Order
}

func NewRepository() *Repository {
	return &Repository{orders: map[domain.OrderID]*domain.Order{}}
}

func (r *Repository) Save(_ context.Context, o *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Lưu bản sao để caller không sửa được dữ liệu đã lưu
	r.orders[o.ID()] = domain.Rehydrate(o.ID(), o.CustomerID(), append([]domain.OrderItem(nil), o.Items()...), o.Status())
	return nil
}

func (r *Repository) FindByID(_ context.Context, id domain.OrderID) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.orders[id]
	if !ok {
		return nil, nil
	}
	return domain.Rehydrate(o.ID(), o.CustomerID(), append([]domain.OrderItem(nil), o.Items()...), o.Status()), nil
}

// PaymentGateway luôn thành công, chỉ ghi log.
type PaymentGateway struct{ Log *slog.Logger }

func (p PaymentGateway) Charge(_ context.Context, customerID uuid.UUID, amount domain.Money) error {
	p.Log.Info("local payment: charged", "customer_id", customerID, "amount", amount.Amount, "currency", amount.Currency)
	return nil
}

// Publisher ghi event ra log.
type Publisher struct{ Log *slog.Logger }

func (p Publisher) Publish(_ context.Context, events []any) error {
	for _, e := range events {
		p.Log.Info("local publisher: event", "event", e)
	}
	return nil
}
